// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package units

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTemperature(t *testing.T) {
	tests := []struct {
		in          string
		want        string
		expectedErr error
	}{
		{
			in:   "0",
			want: "0.0000",
		}, {
			in:   "134F",
			want: "56.6667",
		}, {
			in:   "134 F",
			want: "56.6667",
		}, {
			in:   "134 °F",
			want: "56.6667",
		}, {
			in:   "  134 ° f  ",
			want: "56.6667",
		}, {
			in:   "  134 ° c  ",
			want: "134.0000",
		}, {
			in:   "  -134.1234 ",
			want: "-134.1234",
		}, {
			in:          "134 c°",
			expectedErr: ErrInvalidInput,
		}, {
			in:          "1.3.4",
			expectedErr: ErrInvalidInput,
		},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			assert := assert.New(t)

			got, err := ParseTemperature(tc.in)

			if tc.expectedErr == nil {
				assert.Equal(tc.want, fmt.Sprintf("%.4f", got.C()))
				assert.NoError(err)
				return
			}
			assert.ErrorIs(err, tc.expectedErr)
			assert.Zero(got)
		})
	}
}

func TestTemperature_F(t *testing.T) {
	tests := []struct {
		t    Temperature
		want string
	}{
		{
			t:    Temperature(0),
			want: "32.0000",
		}, {
			t:    Temperature(100),
			want: "212.0000",
		},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tc.want, fmt.Sprintf("%.4f", tc.t.F()))
		})
	}
}

func TestTemperature_ToMug(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{
			in:   "0",
			want: []byte{0x00, 0x00},
		}, {
			in:   "134F",
			want: []byte{0x22, 0x16},
		}, {
			in:   "1F",
			want: []byte{0x00, 0x00},
		}, {
			in:   "1000F",
			want: []byte{0x10, 0x27},
		},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			temp, err := ParseTemperature(tc.in)
			require.NoError(err)

			assert.Equal(tc.want, temp.ToMug())
		})
	}
}

func TestTemperature_FromMug(t *testing.T) {
	tests := []struct {
		out  string
		data []byte
	}{
		{
			out:  "0.0000",
			data: []byte{0x00, 0x00},
		}, {
			out:  "0.0100",
			data: []byte{0x01, 0x00},
		}, {
			out:  "87.2100",
			data: []byte{0x11, 0x22},
		},
	}
	for _, tc := range tests {
		t.Run(tc.out, func(t *testing.T) {
			assert := assert.New(t)

			temp := FromMug(tc.data)
			assert.Equal(tc.out, fmt.Sprintf("%.4f", temp.C()))
		})
	}
}
