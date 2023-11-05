// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/schmidtw/muggo/mug"
	"github.com/schmidtw/muggo/mug/event"
)

func main() {

	m, err := mug.New(
		mug.WithChangeConnectionListener(
			event.ConnectionChangeFunc(
				func(evnt event.ConnectionChange) {
					fmt.Printf("%s, Connected: %v\n", evnt.Address.String(), evnt.Connected)
				},
			)),
	)
	if err != nil {
		panic(err)
	}

	m.Start()

	a := app.New()
	w := a.NewWindow("muggo")

	battery := NewBattery(m)
	battery.Start()
	personalize := NewPersonalize(m, w)
	personalize.Start()
	state := NewState(m)
	state.Start()

	info := container.NewVBox(
		state.Layout(),
		container.NewGridWithColumns(2,
			battery.Layout(),
			personalize.Layout(),
		),
	)

	w.SetContent(info)
	w.ShowAndRun()
}
