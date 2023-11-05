// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"time"

	"github.com/schmidtw/muggo/mug/event"
	bt "tinygo.org/x/bluetooth"
)

func WithAdapter(a *bt.Adapter) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.adapter = a
		return nil
	})
}

func WithAddress(mac string) Option {
	return OptionFunc(func(mug *Mug) error {
		var macAddr bt.MACAddress
		macAddr.Set(mac)
		mug.address = bt.Address{
			MACAddress: macAddr,
		}
		return nil
	})
}

func WithServiceUUIDs(uuids ...string) Option {
	return OptionFunc(func(mug *Mug) error {
		for _, uuid := range uuids {
			if uuid != "" {
				u, err := bt.ParseUUID(uuid)
				if err != nil {
					return err
				}
				mug.serviceUUIDs = append(mug.serviceUUIDs, u)
			}
		}
		return nil
	})
}

type CancelEventListenerFunc func()

func WithChangeConnectionListener(listener event.ConnectionChangeListener, cancel ...*CancelEventListenerFunc) Option {
	return OptionFunc(func(mug *Mug) error {
		cf := mug.changeConnectionListeners.Add(listener)
		if len(cancel) > 0 {
			*cancel[0] = CancelEventListenerFunc(cf)
		}
		return nil
	})
}

func RetryInterval(interval time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.interval = interval
		return nil
	})
}
