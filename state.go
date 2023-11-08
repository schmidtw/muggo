// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/schmidtw/muggo/assets"
	"github.com/schmidtw/muggo/mug"
	"github.com/schmidtw/muggo/mug/event"
	"github.com/schmidtw/muggo/units"
)

const (
	MUG_NONE = iota
	MUG_EMPTY
	MUG_COLD
	MUG_COOL
	MUG_PERFECT
	MUG_WARM
	MUG_HOT
	MUG_COLD_HEATING
	MUG_COOL_HEATING
	MUG_PERFECT_HEATING
)

type State struct {
	m      *mug.Mug
	states map[int]*canvas.Image
	goal   *widget.Entry
	icon   *canvas.Image
	temp   *canvas.Text
	c      *fyne.Container
	rg     *widget.RadioGroup
}

func NewState(m *mug.Mug) *State {
	s := State{
		m: m,
		states: map[int]*canvas.Image{
			MUG_NONE: {
				Resource: fyne.NewStaticResource("mug-disconnected.svg", assets.NoMug),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_EMPTY: {
				Resource: fyne.NewStaticResource("mug-empty.svg", assets.MugEmpty),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COLD: {
				Resource: fyne.NewStaticResource("mug-cold.svg", assets.MugCold),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COOL: {
				Resource: fyne.NewStaticResource("mug-cool.svg", assets.MugCool),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_PERFECT: {
				Resource: fyne.NewStaticResource("mug-perfect.svg", assets.MugPerfect),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_WARM: {
				Resource: fyne.NewStaticResource("mug-warm.svg", assets.MugWarm),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_HOT: {
				Resource: fyne.NewStaticResource("mug-hot.svg", assets.MugHot),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COLD_HEATING: {
				Resource: fyne.NewStaticResource("mug-cold.svg", assets.MugCold),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COOL_HEATING: {
				Resource: fyne.NewStaticResource("mug-cool.svg", assets.MugCool),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_PERFECT_HEATING: {
				Resource: fyne.NewStaticResource("mug-perfect.svg", assets.MugPerfect),
				FillMode: canvas.ImageFillOriginal,
			},
		},
		temp: canvas.NewText(" 78.0 °F", color.White),
		goal: widget.NewEntry(),
	}
	s.rg = widget.NewRadioGroup([]string{"C", "F"}, func(selected string) {
		go func() {
			if selected == "C" {
				s.m.Units(units.Celsius)
			} else {
				s.m.Units(units.Fahrenheit)
			}
		}()
	})

	s.rg.Hidden = false
	s.rg.Horizontal = true
	s.rg.Selected = "C"

	s.temp.Alignment = fyne.TextAlignCenter
	s.temp.TextSize = 48
	s.goal.PlaceHolder = "78.0 °F"
	s.icon = s.states[MUG_NONE]
	s.c = container.NewVBox(
		s.temp,
		s.icon,
		widget.NewForm(
			widget.NewFormItem("Favorite", s.goal),
			widget.NewFormItem("Units", s.rg),
		),
	)
	return &s
}

func (s *State) Start() {
	go func() {
		mugChanges := make(chan mug.MugInfo, 1)
		s.m.AddMugListener(mug.MugListenerFunc(func(info mug.MugInfo) {
			mugChanges <- info
		}))

		conChanges := make(chan event.ConnectionChange, 1)
		s.m.AddConnectionChangeListener(event.ConnectionChangeFunc(func(info event.ConnectionChange) {
			fmt.Printf("connected: %v\n", info.Connected)
			conChanges <- info
		}))

		for {
			var info mug.MugInfo

			s.icon.Refresh()
			s.goal.Refresh()
			s.temp.Refresh()
			s.rg.Refresh()
			if s.c != nil {
				s.c.Objects[0] = s.temp
				s.c.Objects[1] = s.icon
				s.c.Refresh()
			}

			select {
			case info = <-mugChanges:
			case con := <-conChanges:
				if !con.Connected {
					s.icon = s.states[MUG_NONE]
					continue
				}

				info = s.m.All()
			}

			//fmt.Printf("target: %v\n", target)
			//fmt.Printf("current: %v\n", current)
			//fmt.Printf("state: %v\n", state)

			zone := calcTempZone(info.Drink, info.Target)
			switch {
			case info.State == mug.Empty:
				s.icon = s.states[MUG_EMPTY]
			case zone == MUG_COLD:
				s.icon = s.states[MUG_COLD]
				if info.State == mug.Heating {
					s.icon = s.states[MUG_COLD_HEATING]
				}
			case zone == MUG_COOL:
				s.icon = s.states[MUG_COOL]
				if info.State == mug.Heating {
					s.icon = s.states[MUG_COOL_HEATING]
				}
			case zone == MUG_PERFECT:
				s.icon = s.states[MUG_PERFECT]
				if info.State == mug.Heating {
					s.icon = s.states[MUG_PERFECT_HEATING]
				}
			case zone == MUG_WARM:
				s.icon = s.states[MUG_WARM]
			case zone == MUG_HOT:
				s.icon = s.states[MUG_HOT]
			}

			if info.Units == units.Celsius {
				s.rg.Selected = "C"
				s.goal.Text = fmt.Sprintf("%0.01f °C", info.Target.C())
				s.temp.Text = fmt.Sprintf("%0.01f °C", info.Drink.C())
			} else {
				s.rg.Selected = "F"
				s.goal.Text = fmt.Sprintf("%0.01f °F", info.Target.F())
				s.temp.Text = fmt.Sprintf("%0.01f °F", info.Drink.F())
			}
		}
	}()
}

func (s *State) Layout() *fyne.Container {
	return s.c
}

func calcTempZone(current, target units.Temperature) int {
	if -1 < target-current && target-current < 1 {
		return MUG_PERFECT
	}
	if 3 <= current-target {
		return MUG_HOT
	}
	if 0 < current-target {
		return MUG_WARM
	}
	if target-current >= 7 {
		return MUG_COLD
	}
	return MUG_COOL
}
