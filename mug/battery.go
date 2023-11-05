// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"time"

	"github.com/schmidtw/muggo/units"
)

type BatteryInfo struct {
	PercentLeft float64
	Charging    bool
	Temp        units.Temperature
}

func (m *Mug) BatteryInfo() (BatteryInfo, error) {
	data, err := m.io(m, mugApi_BATTERY, 5)
	if err != nil {
		return BatteryInfo{}, err
	}

	return BatteryInfo{
		PercentLeft: float64(data[0]), // 0-100
		Charging:    !(data[1] == 0),  // 0 = not charging, 1 = charging
		Temp:        units.FromMug(data[2:4]),
	}, nil
}

func BatteryTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_BATTERY].ttl = ttl
		return nil
	})
}
