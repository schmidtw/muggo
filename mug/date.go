// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"encoding/binary"
	"time"
)

type DateInfo struct {
	UnixTime int64
	Offset   byte // hours(DST Offset + STD Offset)
}

func (m *Mug) Timestamp(when ...DateInfo) (*DateInfo, error) {
	var write [][]byte

	if len(when) > 0 {
		buf := make([]byte, 0, 5)

		unix := 0xffffffff & when[0].UnixTime
		binary.LittleEndian.PutUint32(buf, uint32(unix))

		buf[4] = when[0].Offset
		write = [][]byte{buf}
	}

	data, err := m.io(m, mugApi_TIME_DATE_ZONE, 5, write...)
	if err != nil {
		return nil, err
	}

	return &DateInfo{
		UnixTime: int64(binary.LittleEndian.Uint32(data)),
		Offset:   data[4],
	}, nil
}

func TimestampTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_TIME_DATE_ZONE].ttl = ttl
		return nil
	})
}
