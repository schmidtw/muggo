// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/schmidtw/muggo/mug/event"
	"github.com/xmidt-org/eventor"
	bt "tinygo.org/x/bluetooth"
)

// Notes:
//
//	Generally the uuid used for the mugs is in the form:
//		fc54%s-236c-4c94-8fa9-944a3e5353fa
const (
	EmberCeramicMugMainServiceUUID   = "fc543622-236c-4c94-8fa9-944a3e5353fa"
	EmberTravelMugMainServiceUUID    = "fc543621-236c-4c94-8fa9-944a3e5353fa"
	EmberTravelMugAltMainServiceUUID = "fc5421a1-236c-4c94-8fa9-944a3e5353fa"
	EmberTravelMugPairServiceUUID    = "fc5421a0-236c-4c94-8fa9-944a3e5353fa"
)

const (
	mugApi_NAME           = 1  // Mug name
	mugApi_DRINK          = 2  // Drink temperature
	mugApi_TARGET         = 3  // Target drink temperature
	mugApi_UNITS          = 4  // Temperature units
	mugApi_LIQUID_LEVEL   = 5  // Liquid level
	mugApi_TIME_DATE_ZONE = 6  // Time date and zone
	mugApi_BATTERY        = 7  // Battery
	mugApi_STATE          = 8  // Liquid state
	mugApi_VOLUME         = 9  // Volume
	mugApi_UNKNOWN_0      = 10 // Unknown
	mugApi_ACCERATION     = 11 // Acceleration
	mugApi_FIRMWARE_INFO  = 12 // Firmware information
	mugApi_ID             = 13 // Mug ID
	mugApi_KEY_0          = 14 // Key 0
	mugApi_KEY_1          = 15 // Key 1
	mugApi_UNKNOWN_1      = 16 // Unknown
	mugApi_UNKNOWN_2      = 17 // Unknown
	mugApi_PUSH_EVENT     = 18 // Push event
	mugApi_UNKNOWN_3      = 19 // Unknown
	mugApi_LED            = 20 // Characteristic LED
	mugApi_LAST           = 21
)

var (
	ErrNotSupported = errors.New("not supported by mug")
	ErrNotConnected = errors.New("not connected to mug")
	ErrInvalidInput = errors.New("invalid input")
)

var (
	Defaults = []Option{
		WithAdapter(bt.DefaultAdapter),
		RetryInterval(5 * time.Second),
		WithServiceUUIDs(
			EmberCeramicMugMainServiceUUID,
			EmberTravelMugMainServiceUUID,
			EmberTravelMugAltMainServiceUUID,
			EmberTravelMugPairServiceUUID,
		),
		// Set the timeouts so they don't all happen at once if possible.
		// Fast moving data
		DrinkTTL(8 * time.Second),
		EmptyTTL(9 * time.Second),
		StateTTL(10 * time.Second),
		BatteryTTL(15 * time.Second),

		// Slow moving data
		NameTTL(24 * time.Hour),
		TargetTTL(24*time.Hour + 1*time.Minute),
		LedTTL(24*time.Hour + 2*time.Minute),
		DeviceInfoTTL(24*time.Hour + 3*time.Minute),
		UnitsTTL(24*time.Hour + 4*time.Minute),
	}
)

type Mug struct {
	m  sync.Mutex
	wg sync.WaitGroup

	adapter  *bt.Adapter
	interval time.Duration

	mugListeners              eventor.Eventor[MugListener]
	changeConnectionListeners eventor.Eventor[event.ConnectionChangeListener]

	address      bt.Address
	serviceUUIDs []bt.UUID

	connShutdown context.CancelFunc

	shutdown context.CancelFunc
	now      func() time.Time

	apis map[int]*cached

	// This makes testing easier because we can mock the device easily.
	io func(m *Mug, api int, length int, write ...[]byte) ([]byte, bool, error)
}

type Option interface {
	apply(*Mug) error
}

type OptionFunc func(*Mug) error

func (f OptionFunc) apply(mug *Mug) error {
	return f(mug)
}

func New(opts ...Option) (*Mug, error) {
	mug := Mug{
		now:  time.Now,
		apis: make(map[int]*cached),
		io:   lockedIO,
	}

	for i := 0; i < mugApi_LAST; i++ {
		mug.apis[i] = &cached{}
	}

	all := append(Defaults, opts...)

	for _, opt := range all {
		if opt != nil {
			err := opt.apply(&mug)
			if err != nil {
				return nil, err
			}
		}
	}

	return &mug, nil
}

func (m *Mug) Start() {
	m.m.Lock()
	defer m.m.Unlock()

	if m.shutdown != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.shutdown = cancel
	go m.run(ctx)
}

func (m *Mug) Stop() {
	m.m.Lock()
	shutdown := m.shutdown
	m.shutdown = nil
	m.m.Unlock()

	if shutdown != nil {
		shutdown()
	}
}

func (m *Mug) AddConnectionChangeListener(listener event.ConnectionChangeListener) {
	m.changeConnectionListeners.Add(listener)
}

func (m *Mug) run(ctx context.Context) {
	m.wg.Add(1)
	defer m.wg.Done()

	for {
		err := m.adapter.Enable()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	disconnected := make(chan struct{}, 1)

	m.adapter.SetConnectHandler(
		func(address bt.Address, connected bool) {
			m.connectHandler(disconnected, address, connected)
		})

	var connected bool

	for {
		if !connected {
			result, err := m.scan(ctx)
			if err == nil {
				err = m.connect(result)
				if err == nil {
					connected = true
				}
			}

			if err != nil {
				fmt.Println(err)
				time.Sleep(m.interval)
				continue
			}
		}

		select {
		case <-ctx.Done():
			return

		case <-disconnected:
			fmt.Println("disconnected in loop")
			connected = false
			m.disconnect()
		}
	}
}

func (m *Mug) connectHandler(notify chan struct{}, address bt.Address, connected bool) {
	m.m.Lock()
	want := m.address
	m.m.Unlock()

	if address != want {
		return
	}

	if !connected {
		notify <- struct{}{}
	}
}

func (m *Mug) scan(ctx context.Context) (*bt.ScanResult, error) {
	type scanResult struct {
		result bt.ScanResult
		err    error
	}

	ch := make(chan scanResult, 1)

	go func() {
		err := m.adapter.Scan(func(adapter *bt.Adapter, r bt.ScanResult) {
			if r.Address == m.address {
				_ = adapter.StopScan()
				ch <- scanResult{result: r}
				return
			}

			for _, service := range m.serviceUUIDs {
				if r.HasServiceUUID(service) {
					_ = adapter.StopScan()
					ch <- scanResult{result: r}
					return
				}
			}
		})
		if err != nil {
			ch <- scanResult{err: err}
		}
	}()

	select {
	case <-ctx.Done():
		m.m.Lock()
		m.address = bt.Address{}
		if m.connShutdown != nil {
			m.connShutdown()
		}
		m.m.Unlock()
		return nil, ctx.Err()
	case result := <-ch:
		if result.err != nil {
			return nil, result.err
		}
		return &result.result, nil
	}
}

func (m *Mug) connect(result *bt.ScanResult) error {
	fmt.Println("found one, connecting")
	address := result.Address
	device, err := m.adapter.Connect(address, bt.ConnectionParams{})
	if err != nil {
		return err
	}

	services, err := device.DiscoverServices(nil)
	if err != nil {
		return err
	}

	m.m.Lock()
	m.address = address
	wait := sync.WaitGroup{}
	for _, service := range services {
		if !m.isWantedService(service.UUID()) {
			continue
		}

		chars, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return err
		}

		for i := range chars {
			char := &chars[i]
			id := uuidToApiId(char.UUID())
			if _, ok := m.apis[id]; !ok {
				m.apis[id] = &cached{}
			}
			m.apis[id].characteristic = char

			// If we are connected, we want to read the mug data without delay.
			wait.Add(1)
			go func() {
				data, err := m.apis[id].read(m.now())
				if err != nil {
					fmt.Printf("id: %d - err: %v\n", id, err)
				} else {
					fmt.Printf("id: %d - data: %x\n", id, data)
				}
				wait.Done()
			}()
		}
	}

	wait.Wait()

	connCtx, cancel := context.WithCancel(context.Background())
	m.connShutdown = cancel
	err = m.startNotifications(connCtx)
	m.m.Unlock()
	m.notifyConnectionChange(true)
	m.dispatch()

	return err
}

func (m *Mug) disconnect() {
	m.m.Lock()
	for k := range m.apis {
		m.apis[k].characteristic = nil
	}
	m.m.Unlock()

	m.notifyConnectionChange(false)
}

func (m *Mug) isWantedService(uuid bt.UUID) bool {
	for _, service := range m.serviceUUIDs {
		if uuid.String() == service.String() {
			return true
		}
	}
	return false
}

func (m *Mug) notifyConnectionChange(connected bool) {
	m.m.Lock()
	address := m.address
	m.m.Unlock()

	cc := event.ConnectionChange{
		Address:   address,
		Connected: connected,
	}
	m.changeConnectionListeners.Visit(func(listener event.ConnectionChangeListener) {
		listener.OnConnectionChange(cc)
	})
}

func uuidToApiId(uuid bt.UUID) int {
	id := uuid.Bytes()
	return (0xff&int(id[13]))<<8 + (0xff & int(id[12]))
}

func lockedIO(m *Mug, api int, length int, write ...[]byte) ([]byte, bool, error) {
	m.m.Lock()
	defer m.m.Unlock()

	prev := m.apis[api].data

	impl, ok := m.apis[api]
	if !ok || impl == nil {
		return nil, false, ErrNotSupported
	}

	if len(write) > 0 {
		_, err := impl.write(write[0])
		if err != nil {
			return nil, false, err
		}
	}

	rv, err := impl.read(m.now())
	if err != nil {
		return nil, false, err
	}

	changed := false
	if len(rv) != len(prev) {
		changed = true
	} else {
		for i := 0; i < len(rv); i++ {
			if rv[i] != prev[i] {
				changed = true
				break
			}
		}
	}

	if length < 1 {
		return rv, changed, err
	}

	if len(rv) != length {
		return nil, false, ErrNotSupported
	}

	return rv, changed, nil
}
