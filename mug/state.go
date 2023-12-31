// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"fmt"
	"time"
)

type State int

const (
	Unknown State = iota
	Empty
	Filling
	Cold
	Cooling
	Heating
	Perfect
	Hot
)

var stateMap = map[byte]State{
	1: Empty,
	2: Filling,
	3: Cold,
	4: Cooling,
	5: Heating,
	6: Perfect,
	7: Hot,
}

var stateStringMap = map[State]string{
	Unknown: "Unknown",
	Empty:   "Empty",
	Filling: "Filling",
	Cold:    "Cold",
	Cooling: "Cooling",
	Heating: "Heating",
	Perfect: "Perfect",
	Hot:     "Hot",
}

func (s State) String() string {
	if rv, ok := stateStringMap[s]; ok {
		return rv
	}

	return fmt.Sprintf("Unknown (%d)", s)
}

func (s State) IsUnknown() bool {
	if _, ok := stateStringMap[s]; ok {
		return false
	}

	return true
}

// State returns the current state of the mug.
func (m *Mug) State() (State, error) {
	data, changed, err := m.io(m, mugApi_STATE, 1)
	if err != nil {
		return Unknown, err
	}

	if changed {
		m.dispatch()
	}

	return stateFromData(data), nil
}

func (m *Mug) stateChanged() {
	m.m.Lock()
	m.apis[mugApi_STATE].expire()
	m.m.Unlock()
	_, _ = m.State()
}

func stateFromData(data []byte) State {
	if state, ok := stateMap[data[0]]; ok {
		return state
	}
	return Unknown
}

// StateTTL sets the TTL for the state of the mug.
func StateTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_STATE].ttl = ttl
		return nil
	})
}
