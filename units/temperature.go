// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package units

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrInvalidInput is returned when the input is not a valid temperature.
	ErrInvalidInput = errors.New("invalid input")
)

type TemperatureUnit string

const (
	Unknown    TemperatureUnit = ""
	Celsius    TemperatureUnit = "C"
	Fahrenheit TemperatureUnit = "F"
)

// Temperature is a temperature in Celsius.
type Temperature float64

// C returns the temperature in Celsius.
func (t Temperature) C() float64 {
	return float64(t)
}

// F returns the temperature in Fahrenheit.
func (t Temperature) F() float64 {
	return float64(t)*9/5 + 32
}

// ParseTemperature parses a string into a Temperature. The string can be in
// either Celsius or Fahrenheit. The string can have a degree symbol or not.
func ParseTemperature(s string) (Temperature, error) {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	isC := true
	if strings.HasSuffix(s, string(Fahrenheit)) {
		isC = false
	}
	if strings.HasSuffix(s, "°") {
		return 0, ErrInvalidInput
	}
	s = strings.Replace(s, "°", "", -1)
	s = strings.TrimSuffix(s, string(Fahrenheit))
	s = strings.TrimSuffix(s, string(Celsius))
	s = strings.TrimSpace(s)

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Join(ErrInvalidInput, err)
	}

	if isC {
		return Temperature(num), nil
	}

	return Temperature((num - 32) * 5 / 9), nil
}

// ToMug returns the byte representation of the temperature in increments of
// 0.01 C little endian.
func (t Temperature) ToMug() []byte {
	temp := int(t * 100)
	// It can't be colder than frozen or hotter than boiling.
	if temp < 0 {
		temp = 0
	}
	if temp > 10000 {
		temp = 10000
	}

	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(temp))

	return buf
}

// FromMug converts the byte representation of the temperature in increments of
// 0.01 C little endian to a Temperature.
func (t *Temperature) FromMug(data []byte) {
	*t = Temperature(binary.LittleEndian.Uint16(data)) * 0.01
}

// FromMug converts the byte representation of the temperature in increments of
// 0.01 C little endian to a Temperature.
func FromMug(data []byte) Temperature {
	var t Temperature
	t.FromMug(data)
	return t
}
