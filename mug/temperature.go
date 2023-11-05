// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"context"
	"time"

	"github.com/schmidtw/muggo/units"
)

// Current returns the current temperature of the mug in celcius.
func (m *Mug) Current(ctx context.Context) (units.Temperature, error) {
	data, err := m.io(m, mugApi_TEMP, 2)
	if err != nil {
		return 0.0, err
	}

	return units.FromMug(data), nil
}

// CurrentTTL sets the TTL for the temperature of the drink in the mug.
func CurrentTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_TEMP].ttl = ttl
		return nil
	})
}
