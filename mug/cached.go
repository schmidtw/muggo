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

func (c *cached) returnCached(now time.Time) bool {
	if c.fetched.IsZero() || c.fetched.Add(c.ttl).Before(now) {
		return false
	}
	return true
}

func (c *cached) expire() {
	c.fetched = time.Time{}
}

func (c *cached) read(now time.Time) ([]byte, error) {
	if c.characteristic == nil {
		return nil, ErrNotConnected
	}

	// If now is zero, we are forcing a fetch.
	if !now.IsZero() && c.returnCached(now) {
		return c.data, nil
	}

	max, err := c.characteristic.GetMTU()
	if err != nil {
		return nil, err
	}

	data := make([]byte, max)
	len, err := c.characteristic.Read(data)
	if err != nil {
		return nil, err
	}

	c.data = data[:len]
	c.fetched = now

	return c.data, nil
}

func (c *cached) write(data []byte) (int, error) {
	if c.characteristic == nil {
		return 0, ErrNotConnected
	}

	l, err := c.characteristic.WriteWithoutResponse(data)
	if err != nil {
		c.fetched = time.Time{}
	}
	return l, err
}
