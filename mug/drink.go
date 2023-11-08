// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"time"

	"github.com/schmidtw/muggo/units"
)

// Drink returns the current temperature of the drink in celsius.
func (m *Mug) Drink() (units.Temperature, error) {
	data, changed, err := m.io(m, mugApi_DRINK, 2)
	if err != nil {
		return 0.0, err
	}

	if changed {
		m.dispatch()
	}

	return drinkFromData(data), nil
}

func (m *Mug) drinkChanged() {
	m.m.Lock()
	m.apis[mugApi_DRINK].expire()
	m.m.Unlock()
	_, _ = m.Drink()
}

func drinkFromData(data []byte) units.Temperature {
	return units.FromMug(data)
}

// DrinkTTL sets the TTL for the temperature of the drink in the mug.
func DrinkTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_DRINK].ttl = ttl
		return nil
	})
}
