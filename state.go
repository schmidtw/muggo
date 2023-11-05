// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/schmidtw/muggo/assets"
	"github.com/schmidtw/muggo/mug"
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
			MUG_NONE: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-disconnected.svg", assets.NoMug),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_EMPTY: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-empty.svg", assets.MugEmpty),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COLD: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-cold.svg", assets.MugCold),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COOL: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-cool.svg", assets.MugCool),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_PERFECT: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-perfect.svg", assets.MugPerfect),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_WARM: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-warm.svg", assets.MugWarm),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_HOT: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-hot.svg", assets.MugHot),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COLD_HEATING: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-cold.svg", assets.MugCold),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_COOL_HEATING: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-cool.svg", assets.MugCool),
				FillMode: canvas.ImageFillOriginal,
			},
			MUG_PERFECT_HEATING: &canvas.Image{
				Resource: fyne.NewStaticResource("mug-perfect.svg", assets.MugPerfect),
				FillMode: canvas.ImageFillOriginal,
			},
		},
		temp: canvas.NewText(" 78.0 °F", color.White),
		goal: widget.NewEntry(),
		rg:   widget.NewRadioGroup([]string{"C", "F"}, func(string) {}),
	}

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
		var connected bool
		var unitsRead bool

		current := units.Temperature(25.0)
		target := units.Temperature(54.0)
		state := mug.Empty
		unit := units.Celsius

		for {

			//fmt.Printf("target: %v\n", target)
			//fmt.Printf("current: %v\n", current)
			//fmt.Printf("state: %v\n", state)

			zone := calcTempZone(current, target)
			switch {
			case !connected:
				s.icon = s.states[MUG_NONE]
			case state == mug.Empty:
				s.icon = s.states[MUG_EMPTY]
			case zone == MUG_COLD:
				s.icon = s.states[MUG_COLD]
				if state == mug.Heating {
					s.icon = s.states[MUG_COLD_HEATING]
				}
			case zone == MUG_COOL:
				s.icon = s.states[MUG_COOL]
				if state == mug.Heating {
					s.icon = s.states[MUG_COOL_HEATING]
				}
			case zone == MUG_PERFECT:
				s.icon = s.states[MUG_PERFECT]
				if state == mug.Heating {
					s.icon = s.states[MUG_PERFECT_HEATING]
				}
			case zone == MUG_WARM:
				s.icon = s.states[MUG_WARM]
			case zone == MUG_HOT:
				s.icon = s.states[MUG_HOT]
			}

			if s.rg.Selected == "C" {
				s.goal.Text = fmt.Sprintf("%0.01f °C", target.C())
				s.temp.Text = fmt.Sprintf("%0.01f °C", current.C())
			} else {
				s.goal.Text = fmt.Sprintf("%0.01f °F", target.F())
				s.temp.Text = fmt.Sprintf("%0.01f °F", current.F())
			}
			s.icon.Refresh()
			s.goal.Refresh()
			s.temp.Refresh()
			if s.c != nil {
				s.c.Objects[0] = s.temp
				s.c.Objects[1] = s.icon
				s.c.Refresh()
			}

			time.Sleep(1 * time.Second)

			var err error

			target, _ = s.m.Target()
			current, _ = s.m.Current(context.Background())
			unit, err = s.m.Units()
			if err == nil {
				if unitsRead {
					out := units.Celsius
					if s.rg.Selected == "F" {
						out = units.Fahrenheit
					}
					_, _ = s.m.Units(out)
				} else {
					unitsRead = true
					s.rg.Selected = "C"
					if unit == units.Fahrenheit {
						s.rg.Selected = "F"
					}
				}
			}
			state, err = s.m.State()

			connected = true
			if err != nil {
				connected = false
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
