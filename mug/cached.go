// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"time"

	bt "tinygo.org/x/bluetooth"
)

type cached struct {
	characteristic *bt.DeviceCharacteristic
	data           []byte
	fetched        time.Time
	ttl            time.Duration
}

func (m *cached) returnCached(now time.Time) bool {
	if m.fetched.IsZero() || m.fetched.Add(m.ttl).Before(now) {
		return false
	}
	return true
}

func (m *cached) read(now time.Time) ([]byte, error) {
	if m.characteristic == nil {
		return nil, ErrNotConnected
	}

	// If now is zero, we are forcing a fetch.
	if !now.IsZero() && m.returnCached(now) {
		return m.data, nil
	}

	max, err := m.characteristic.GetMTU()
	if err != nil {
		return nil, err
	}

	data := make([]byte, max)
	len, err := m.characteristic.Read(data)
	if err != nil {
		return nil, err
	}

	m.data = data[:len]
	m.fetched = now

	return m.data, nil
}

func (m *cached) write(data []byte) (int, error) {
	l, err := m.characteristic.WriteWithoutResponse(data)
	if err != nil {
		m.fetched = time.Time{}
	}
	return l, err
}
