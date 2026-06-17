// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// async kører work() i en goroutine (så UI'en ikke fryser under et CLI-kald) og
// leverer resultatet tilbage på UI-tråden via fyne.Do, hvor done() trygt kan
// røre widgets.
func (u *ui) async(work func() any, done func(any)) {
	go func() {
		v := work()
		fyne.Do(func() { done(v) })
	}()
}

// newActivityLog laver det read-only tekstfelt der viser CLI-output.
func newActivityLog() *widget.Entry {
	e := widget.NewMultiLineEntry()
	e.Wrapping = fyne.TextWrapWord
	e.Disable() // read-only, men stadig scrollbar/markérbar
	return e
}

// log tilføjer en linje (med tidsstempel) til aktivitetsloggen.
func (u *ui) log(text string) {
	text = strings.TrimRight(text, "\n")
	if text == "" {
		return
	}
	stamp := time.Now().Format("15:04:05")
	cur := u.activity.Text
	if cur != "" {
		cur += "\n"
	}
	u.activity.SetText(cur + "[" + stamp + "] " + text)
	u.activity.CursorRow = strings.Count(u.activity.Text, "\n")
	u.activity.Refresh()
}

// centeredSpinner viser en uendelig progress-bar med en label — bruges mens
// vi venter på et CLI-kald.
func centeredSpinner(label string) fyne.CanvasObject {
	bar := widget.NewProgressBarInfinite()
	box := container.NewVBox(widget.NewLabel(label), bar)
	return container.New(layout.NewCenterLayout(), box)
}

// errSimple pakker en streng som en error til dialog.ShowError.
func errSimple(msg string) error { return errors.New(msg) }

// uriToPath konverterer en Fyne-URI fra en fildialog til en OS-sti. På Windows
// kan Fyne give "/C:/sti" — den ledende skråstreg fjernes så Go's exec/os kan
// bruge stien.
func uriToPath(u fyne.URI) string {
	p := u.Path()
	if len(p) >= 3 && p[0] == '/' && p[2] == ':' {
		p = p[1:]
	}
	return filepath.FromSlash(p)
}
