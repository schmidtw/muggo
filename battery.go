// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/schmidtw/muggo/assets"
	"github.com/schmidtw/muggo/mug"
	"github.com/schmidtw/muggo/mug/event"
)

type Battery struct {
	m           *mug.Mug
	info        mug.BatteryInfo
	discharging []*canvas.Image
	charging    []*canvas.Image
	icon        *canvas.Image
	text        *canvas.Text
	c           *fyne.Container
}

func NewBattery(m *mug.Mug) *Battery {
	b := Battery{
		m: m,
		discharging: []*canvas.Image{
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-alert-0.svg", assets.BatteryAlert0),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-alert-1.svg", assets.BatteryAlert1),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-normal-2.svg", assets.BatteryNormal2),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-normal-3.svg", assets.BatteryNormal3),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-normal-4.svg", assets.BatteryNormal4),
			),
		},
		charging: []*canvas.Image{
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-charging-0.svg", assets.BatteryCharging0),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-charging-1.svg", assets.BatteryCharging1),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-charging-2.svg", assets.BatteryCharging2),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-charging-3.svg", assets.BatteryCharging3),
			),
			canvas.NewImageFromResource(
				fyne.NewStaticResource("battery-charging-4.svg", assets.BatteryCharging4),
			),
		},
		text: canvas.NewText("0%", color.White),
	}

	b.icon = b.discharging[0]

	b.update(mug.BatteryInfo{})

	return &b
}

func (b *Battery) Start() {
	go func() {
		mugChanges := make(chan mug.BatteryInfo, 1)
		b.m.AddMugListener(mug.MugListenerFunc(func(info mug.MugInfo) {
			mugChanges <- info.Battery
		}))

		conChanges := make(chan event.ConnectionChange, 1)
		b.m.AddConnectionChangeListener(event.ConnectionChangeFunc(func(info event.ConnectionChange) {
			conChanges <- info
		}))

		for {
			var info mug.BatteryInfo

			select {
			case info = <-mugChanges:
			case con := <-conChanges:
				if !con.Connected {
					continue
				}

				all := b.m.All()
				info = all.Battery
			}

			b.update(info)
		}
	}()
}

func (b *Battery) Layout() *fyne.Container {
	b.c = container.New(layout.NewVBoxLayout(), b.icon, b.text)
	return b.c
}

// SetLevel sets the battery level to the given value (0.0 to 1.0)
func (b *Battery) update(bi mug.BatteryInfo) {
	b.info = bi

	choices := b.discharging
	if b.info.Charging {
		choices = b.charging
	}

	which := len(choices) - 1
	for i := 0; i < len(choices); i++ {
		area := 100.0 / float64(len(choices))
		if b.info.PercentLeft < area*float64(i+1) {
			which = i
			break
		}
	}

	b.icon = choices[which]
	b.icon.FillMode = canvas.ImageFillOriginal

	b.text.Text = fmt.Sprintf("%.0f%%", b.info.PercentLeft)
	b.text.Alignment = fyne.TextAlignCenter
	b.icon.Refresh()
	b.text.Refresh()
	if b.c != nil {
		b.c.Objects[0] = b.icon
		b.c.Refresh()
	}
}
