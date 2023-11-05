// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"time"
)

// Name returns the name of the mug as a string.  If the optional name
// parameter is provided, it will be used to set the name of the mug first.
func (m *Mug) Name(name ...string) (string, error) {
	var write [][]byte
	if len(name) > 0 {
		write = [][]byte{[]byte(name[0])}
	}

	data, err := m.io(m, mugApi_NAME, 0, write...)

	return string(data), err
}

// NameTTL sets the TTL for the name of the mug.
func NameTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_NAME].ttl = ttl
		return nil
	})
}
