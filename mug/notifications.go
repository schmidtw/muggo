// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"context"
	"fmt"
	"time"
)

const (
	NOTIFY_BATTERY              = 1
	NOTIFY_CHARGHING            = 2
	NOTIFY_DISCHARGING          = 3
	NOTIFY_TARGET_CHANGED       = 4
	NOTIFY_DRINK_CHANGED        = 5
	NOTIFY_UNSURE               = 6 // Not sure what this is
	NOTIFY_LIQUID_LEVEL_CHANGED = 7
	NOTIFY_STATE_CHANGED        = 8
)

func (m *Mug) startNotifications(ctx context.Context) error {
	err := m.apis[mugApi_PUSH_EVENT].characteristic.EnableNotifications(
		func(buf []byte) {

			switch buf[0] {
			case NOTIFY_BATTERY:
				go m.refreshbattery()

			case NOTIFY_CHARGHING:
				go m.charging()

			case NOTIFY_DISCHARGING:
				go m.discharging()

			case NOTIFY_TARGET_CHANGED:
				fmt.Println("notify: target changed")
				go m.targetChanged()

			case NOTIFY_DRINK_CHANGED:
				go m.drinkChanged()

			case NOTIFY_LIQUID_LEVEL_CHANGED:
				go m.emptyChanged()

			case NOTIFY_STATE_CHANGED:
				fmt.Println("notify: state changed")
				go m.stateChanged()
			default:
				fmt.Println("unknown push event:", buf)
			}
		},
	)

	go m.handleMissingNotifications(ctx)

	return err
}

func (m *Mug) handleMissingNotifications(ctx context.Context) {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = m.Drink()
			_, _ = m.Target()
			_, _ = m.BatteryInfo()
			_, _ = m.IsEmpty()
			_, _ = m.State()
			_, _ = m.Units()
		}
	}
}
