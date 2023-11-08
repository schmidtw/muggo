// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"fmt"
	"time"

	"github.com/schmidtw/muggo/units"
)

// Units returns the current units of the mug.  If a unit is provided, the mug
// will be set to that unit.
// This function only appears to be useful if the mug has a display.
func (m *Mug) Units(unit ...units.TemperatureUnit) (units.TemperatureUnit, error) {
	var write [][]byte
	if len(unit) > 0 {
		t := []byte{0}
		switch unit[0] {
		case units.Fahrenheit:
			t[0] = 1
		case units.Celsius:
			t[0] = 0
		default:
			return units.Unknown, ErrInvalidInput
		}
		write = [][]byte{t}
	}

	data, _, err := m.io(m, mugApi_UNITS, 2, write...)
	if err != nil {
		return units.Unknown, err
	}

	rv := unitsFromData(data)
	if rv != units.Unknown {
		return rv, nil
	}

	return rv, fmt.Errorf("%w unknown response: 0x%02x", ErrNotSupported, data[0])
}

func unitsFromData(data []byte) units.TemperatureUnit {
	switch data[0] {
	case 0:
		return units.Celsius
	case 1:
		return units.Fahrenheit
	default:
	}

	return units.Unknown
}

// UnitsTTL sets the TTL for the units command, so it can be cached.
func UnitsTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_UNITS].ttl = ttl
		return nil
	})
}
