// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package mug

import (
	"encoding/binary"
	"strings"
	"time"
)

type DeviceInfo struct {
	FirmwareVersion   uint16
	HardwareVersion   uint16
	BootloaderVersion uint16
	SerialNumber      string
}

func (m *Mug) DeviceInfo() (*DeviceInfo, error) {
	data, _, err := m.io(m, mugApi_FIRMWARE_INFO, 0)
	if err != nil {
		return nil, err
	}

	serial, _, err := m.io(m, mugApi_ID, 0)
	if err != nil {
		return nil, err
	}

	di := deviceInfoFromData(data, serial)
	if di != nil {
		return di, nil
	}

	return nil, ErrNotSupported
}

func deviceInfoFromData(data, serial []byte) *DeviceInfo {
	var serialNumber string
	if len(serial) > 6 {
		// It's not clear what the first 6 bytes are, but they are not
		// part of the serial number.
		serialNumber = strings.Replace(string(serial[6:]), "-", "", -1)
	}

	if len(data) == 4 {
		return &DeviceInfo{
			FirmwareVersion: binary.LittleEndian.Uint16(data[0:2]),
			HardwareVersion: binary.LittleEndian.Uint16(data[2:4]),
			SerialNumber:    serialNumber,
		}
	}

	if len(data) == 6 {
		return &DeviceInfo{
			FirmwareVersion:   binary.LittleEndian.Uint16(data[0:2]),
			HardwareVersion:   binary.LittleEndian.Uint16(data[2:4]),
			BootloaderVersion: binary.LittleEndian.Uint16(data[4:6]),
			SerialNumber:      serialNumber,
		}
	}

	return nil
}

func DeviceInfoTTL(ttl time.Duration) Option {
	return OptionFunc(func(mug *Mug) error {
		mug.apis[mugApi_FIRMWARE_INFO].ttl = ttl
		mug.apis[mugApi_ID].ttl = ttl
		return nil
	})
}
