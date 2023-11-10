// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package assets

import (
	_ "embed"
	_ "image/png"
)

//go:embed battery-alert-0.svg
var BatteryAlert0 []byte

//go:embed battery-alert-1.svg
var BatteryAlert1 []byte

//go:embed battery-normal-2.svg
var BatteryNormal2 []byte

//go:embed battery-normal-3.svg
var BatteryNormal3 []byte

//go:embed battery-normal-4.svg
var BatteryNormal4 []byte

//go:embed battery-charging-0.svg
var BatteryCharging0 []byte

//go:embed battery-charging-1.svg
var BatteryCharging1 []byte

//go:embed battery-charging-2.svg
var BatteryCharging2 []byte

//go:embed battery-charging-3.svg
var BatteryCharging3 []byte

//go:embed battery-charging-4.svg
var BatteryCharging4 []byte

//go:embed mug-disconnected.svg
var NoMug []byte

//go:embed mug-empty.svg
var MugEmpty []byte

//go:embed mug-cold.svg
var MugCold []byte

//go:embed mug-cool.svg
var MugCool []byte

//go:embed mug-perfect.svg
var MugPerfect []byte

//go:embed mug-warm.svg
var MugWarm []byte

//go:embed mug-hot.svg
var MugHot []byte

//go:embed mug-cold-heating.svg
var MugColdHeating []byte

//go:embed mug-cool-heating.svg
var MugCoolHeating []byte

//go:embed mug-perfect-heating.svg
var MugPerfectHeating []byte
