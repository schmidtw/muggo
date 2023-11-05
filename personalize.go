// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
	"github.com/schmidtw/muggo/mug"
)

type Personal struct {
	m          *mug.Mug
	led        *color.NRGBA
	name       string
	nameWidget *widget.Entry
	ledWidget  colorpicker.PickerOpenWidget
	c          *fyne.Container
}

func NewPersonalize(m *mug.Mug, w fyne.Window) *Personal {
	p := Personal{
		m:          m,
		nameWidget: widget.NewEntry(),
	}

	size := p.nameWidget.MinSize()
	p.nameWidget.SetPlaceHolder("________________")
	p.nameWidget.Resize(
		fyne.Size{
			Height: size.Height,
			Width:  200,
		})
	p.ledWidget = colorpicker.NewColorSelectModalRect(w,
		fyne.Size{
			Height: size.Height,
			Width:  size.Height,
		},
		color.Black)

	return &p
}

func (p *Personal) Start() {
	go func() {
		for {
			var err error
			time.Sleep(1 * time.Second)
			p.name, err = p.m.Name()
			if err == nil {
				p.nameWidget.SetText(p.name)
				p.nameWidget.Refresh()
			}
			p.led, err = p.m.Led()
			if err == nil {
				p.ledWidget.SetColor(p.led)
				p.ledWidget.Refresh()
			}
		}
	}()
}

func (p *Personal) Layout() *fyne.Container {
	p.c = container.NewPadded(p.ledWidget)
	return p.c
}
