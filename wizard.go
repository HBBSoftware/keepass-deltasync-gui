// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// showWizard er onboarding-guiden for en bruger der endnu ikke har tilmeldt sin
// enhed. Den fører gennem to trin: (1) tilmeld enhed med server-adresse + token,
// (2) tilføj den første lokale database. Begge trin udfører arbejdet ved at
// kalde CLI'en.

// showWizard viser velkomstskærmen.
func (u *ui) showWizard() {
	title := widget.NewLabelWithStyle(L.WizardWelcomeTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	body := widget.NewLabel(L.WizardWelcomeBody)
	body.Wrapping = fyne.TextWrapWord

	start := widget.NewButton(L.WizardStart, func() { u.showEnrollStep() })
	start.Importance = widget.HighImportance

	card := container.NewVBox(
		title,
		widget.NewSeparator(),
		body,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), start, layout.NewSpacer()),
	)
	u.win.SetContent(container.NewPadded(card))
}

// showEnrollStep er trin 1: tilmeld enheden hos serveren.
func (u *ui) showEnrollStep() {
	server := widget.NewEntry()
	server.SetPlaceHolder(L.ServerURLHint)
	token := widget.NewEntry()
	token.SetPlaceHolder(L.EnrollTokenHint)
	device := widget.NewEntry()
	device.SetPlaceHolder(L.DeviceNameHint)

	form := widget.NewForm(
		widget.NewFormItem(L.ServerURL, server),
		widget.NewFormItem(L.EnrollToken, token),
		widget.NewFormItem(L.DeviceName, device),
	)

	submit := widget.NewButton(L.EnrollButton, func() {
		if server.Text == "" || token.Text == "" {
			dialog.ShowError(errSimple(L.ServerURL+" / "+L.EnrollToken), u.win)
			return
		}
		u.win.SetContent(centeredSpinner(L.Working))
		u.async(func() any {
			ctx, cancel := withTimeout(60 * time.Second)
			defer cancel()
			return u.c.enroll(ctx, server.Text, device.Text, token.Text)
		}, func(v any) {
			r := v.(result)
			u.log(r.Combined())
			if r.Err != nil {
				u.showEnrollStep()
				dialog.ShowError(errSimple(describeErr(r)), u.win)
				return
			}
			u.showAddDBStep()
		})
	})
	submit.Importance = widget.HighImportance

	header := widget.NewLabelWithStyle(L.WizardStepEnroll, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	back := widget.NewButton(L.Back, func() { u.showWizard() })

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		form,
		layout.NewSpacer(),
		container.NewHBox(back, layout.NewSpacer(), submit),
	)
	u.win.SetContent(container.NewPadded(content))
}

// showAddDBStep er trin 2: registrér den første lokale .kdbx-fil.
func (u *ui) showAddDBStep() {
	name := widget.NewEntry()
	name.SetPlaceHolder(L.DBNameHint)
	path := widget.NewEntry()
	path.SetPlaceHolder(L.KdbxFileHint)

	browse := widget.NewButton(L.Browse, func() {
		dialog.ShowFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			defer rc.Close()
			path.SetText(uriToPath(rc.URI()))
		}, u.win)
	})
	pathRow := container.NewBorder(nil, nil, nil, browse, path)

	body := widget.NewLabel(L.WizardAddDBBody)
	body.Wrapping = fyne.TextWrapWord

	form := widget.NewForm(
		widget.NewFormItem(L.DBName, name),
		widget.NewFormItem(L.KdbxFile, pathRow),
	)

	create := widget.NewButton(L.CreateDBButton, func() {
		if name.Text == "" || path.Text == "" {
			dialog.ShowError(errSimple(L.DBName+" / "+L.KdbxFile), u.win)
			return
		}
		u.win.SetContent(centeredSpinner(L.Working))
		u.async(func() any {
			ctx, cancel := withTimeout(60 * time.Second)
			defer cancel()
			return u.c.initDB(ctx, name.Text, path.Text)
		}, func(v any) {
			r := v.(result)
			u.log(r.Combined())
			if r.Err != nil {
				u.showAddDBStep()
				dialog.ShowError(errSimple(describeErr(r)), u.win)
				return
			}
			u.showWizardDone()
		})
	})
	create.Importance = widget.HighImportance

	skip := widget.NewButton(L.SkipForNow, func() { u.showDashboard() })

	header := widget.NewLabelWithStyle(L.WizardStepDB, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		widget.NewLabelWithStyle(L.WizardAddDB, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		body,
		form,
		layout.NewSpacer(),
		container.NewHBox(skip, layout.NewSpacer(), create),
	)
	u.win.SetContent(container.NewPadded(content))
}

// showWizardDone er afslutningsskærmen.
func (u *ui) showWizardDone() {
	title := widget.NewLabelWithStyle(L.WizardDone, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	body := widget.NewLabel(L.WizardDoneBody)
	body.Wrapping = fyne.TextWrapWord
	finish := widget.NewButton(L.Finish, func() { u.showDashboard() })
	finish.Importance = widget.HighImportance

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		body,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), finish, layout.NewSpacer()),
	)
	u.win.SetContent(container.NewPadded(content))
}
