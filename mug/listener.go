// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"image/color"

	"github.com/schmidtw/muggo/units"
)

type MugInfo struct {
	Name       string
	Drink      units.Temperature
	Target     units.Temperature
	Battery    BatteryInfo
	Empty      bool
	LED        color.NRGBA
	DeviceInfo DeviceInfo
	State      State
	Units      units.TemperatureUnit
}

type MugListener interface {
	MugInfo(MugInfo)
}

type CancelFunc func()

func (m *Mug) AddMugListener(l MugListener) CancelFunc {
	cancel := m.mugListeners.Add(l)
	return CancelFunc(cancel)
}

type MugListenerFunc func(MugInfo)

func (f MugListenerFunc) MugInfo(m MugInfo) {
	f(m)
}

func (m *Mug) All() MugInfo {
	m.m.Lock()
	defer m.m.Unlock()

	di := deviceInfoFromData(
		m.apis[mugApi_FIRMWARE_INFO].data,
		m.apis[mugApi_ID].data,
	)
	if di == nil {
		di = &DeviceInfo{}
	}

	return MugInfo{
		Name:       nameFromData(m.apis[mugApi_NAME].data),
		Drink:      drinkFromData(m.apis[mugApi_DRINK].data),
		Target:     targetFromData(m.apis[mugApi_TARGET].data),
		Battery:    batteryInfoFromData(m.apis[mugApi_BATTERY].data),
		Empty:      emptyFromData(m.apis[mugApi_LIQUID_LEVEL].data),
		LED:        ledFromData(m.apis[mugApi_LED].data),
		DeviceInfo: *di,
		State:      stateFromData(m.apis[mugApi_STATE].data),
		Units:      unitsFromData(m.apis[mugApi_UNITS].data),
	}
}

func (m *Mug) dispatch() {
	mi := m.All()

	m.mugListeners.Visit(func(l MugListener) {
		l.MugInfo(mi)
	})
}
