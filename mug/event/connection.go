// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package event

import bt "tinygo.org/x/bluetooth"

type ConnectionChange struct {
	Address   bt.Address
	Connected bool
}

type ConnectionChangeListener interface {
	OnConnectionChange(ConnectionChange)
}

// StatusEventFunc is a convenience type for implementing the StatusEventListener
// interface with a function.
type ConnectionChangeFunc func(ConnectionChange)

func (f ConnectionChangeFunc) OnConnectionChange(c ConnectionChange) {
	f(c)
}
