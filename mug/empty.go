// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"time"
)

// Empty returns if the mug is empty.
func (m *Mug) IsEmpty() (bool, error) {
	data, changed, err := m.io(m, mugApi_LIQUID_LEVEL, 1)
	if err != nil {
		return false, err
	}

	if changed {
		m.dispatch()
	}

	return emptyFromData(data), nil
}

func (m *Mug) emptyChanged() {
	m.m.Lock()
	m.apis[mugApi_LIQUID_LEVEL].expire()
	m.m.Unlock()
	_, _ = m.IsEmpty()
}

func emptyFromData(data []byte) bool {
	return (data[0] == 0x00)
}

// DrinkTempTTL sets the TTL for the temperature of the drink in the mug.
func EmptyTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_LIQUID_LEVEL].ttl = ttl
		return nil
	})
}
