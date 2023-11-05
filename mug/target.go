// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"time"

	"github.com/schmidtw/muggo/units"
)

// Target returns the target temperature of the mug.  If a temperature is
// provided, the mug will be set to that temperature.
func (m *Mug) Target(temp ...units.Temperature) (units.Temperature, error) {
	var write [][]byte
	if len(temp) > 0 {
		write = [][]byte{temp[0].ToMug()}
	}

	data, err := m.io(m, mugApi_TARGET, 2, write...)
	if err != nil {
		return 0, err
	}

	return units.FromMug(data), nil
}

// TargetTTL sets the TTL for the target temperature of the mug.
func TargetTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_TARGET].ttl = ttl
		return nil
	})
}
