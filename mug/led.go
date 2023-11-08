// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"image/color"
	"time"
)

func (m *Mug) Led(rgba ...color.NRGBA) (*color.NRGBA, error) {
	var write [][]byte

	if len(rgba) > 0 {
		write = [][]byte{[]byte{rgba[0].R, rgba[0].G, rgba[0].B, rgba[0].A}}
	}

	data, _, err := m.io(m, mugApi_LED, 4, write...)
	if err != nil {
		return nil, err
	}

	rv := ledFromData(data)

	return &rv, nil
}

func ledFromData(data []byte) color.NRGBA {
	return color.NRGBA{
		R: data[0],
		G: data[1],
		B: data[2],
		A: data[3],
	}
}

func LedTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_LED].ttl = ttl
		return nil
	})
}
