// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// showDashboard er hovedskærmen efter tilmelding: en faneblads-visning med
// databaser, aktivitetslog og indstillinger.
func (u *ui) showDashboard() {
	debugf("showDashboard: begin")
	u.selDB = -1
	u.selDev = -1
	u.statusLb = widget.NewLabel("")
	u.statusLb.Wrapping = fyne.TextWrapWord

	dbTab := u.databasesTab()
	debugf("showDashboard: databasesTab built")
	devTab := u.devicesTab()
	debugf("showDashboard: devicesTab built")
	actTab := u.activityTab()
	debugf("showDashboard: activityTab built")
	setTab := u.settingsTab()
	debugf("showDashboard: settingsTab built")

	tabs := container.NewAppTabs(
		container.NewTabItem(L.TabDatabases, dbTab),
		container.NewTabItem(L.TabDevices, devTab),
		container.NewTabItem(L.TabActivity, actTab),
		container.NewTabItem(L.TabSettings, setTab),
	)
	u.win.SetContent(tabs)
	debugf("showDashboard: content set")

	u.refreshStatus()
	u.refreshDatabases()
	u.refreshDevices()
	debugf("showDashboard: refresh kicked off")
}

// dbColumns er kolonnerne i gridden — bevidst de SAMME som `databases`-kommandoen
// i terminalen, så man genkender visningen. width er pixelbredden.
var dbColumns = []gridColumn{
	{func() string { return L.ColStatus }, 120},
	{func() string { return L.ColName }, 140},
	{func() string { return L.ColID }, 290},
	{func() string { return L.ColCreated }, 215},
	{func() string { return L.ColPath }, 340},
}

// cellText giver teksten for en celle i gridden.
func (u *ui) cellText(row, col int) string {
	if row < 0 || row >= len(u.dbList) {
		return ""
	}
	db := u.dbList[row]
	switch col {
	case 0:
		if db.Bound {
			return L.BoundLocally
		}
		return L.OnServerOnly
	case 1:
		return db.Name
	case 2:
		return db.ID
	case 3:
		return db.Created
	case 4:
		return db.LocalPath
	}
	return ""
}

// databasesTab bygger fanen: en grid med de tilkoblede databaser (samme kolonner
// som terminalen) plus handlingsknapper. Man vælger en række i gridden og
// synkroniserer den valgte — eller synkroniserer alle bundne på én gang.
func (u *ui) databasesTab() fyne.CanvasObject {
	u.dbInfo = widget.NewLabel("")

	var gridBody fyne.CanvasObject
	u.dbGrid, gridBody = newGridList(
		dbColumns,
		func() int { return len(u.dbList) },
		u.cellText,
		func(row int) { u.selDB = row; u.loadMembers() },
		func() { u.selDB = -1; u.clearMembers() },
	)

	refresh := widget.NewButton(L.Refresh, func() { u.refreshDatabases() })
	add := widget.NewButton(L.AddDatabase, func() { u.showAddDBDialog() })
	add.Importance = widget.HighImportance
	forget := widget.NewButton(L.ForgetDatabase, func() { u.forgetSelected() })
	syncSel := widget.NewButton(L.SyncSelected, func() { u.syncSelected() })
	syncAll := widget.NewButton(L.SyncAll, func() { u.syncAll() })

	toolbar := container.NewHBox(add, forget, syncSel, syncAll, layout.NewSpacer(), refresh)
	top := container.NewVBox(toolbar, u.dbInfo)

	// Øvre del: database-gridden. Nedre del: detalje-panel for den valgte
	// database (medlemmer/enheder koblet til den).
	gridPane := container.NewBorder(top, nil, nil, nil, gridBody)
	split := container.NewVSplit(gridPane, u.membersPane())
	split.SetOffset(0.6)
	return split
}

// memberColumns er kolonnerne i medlems-gridden — samme som `shares`-kommandoen.
var memberColumns = []gridColumn{
	{func() string { return L.ColRole }, 110},
	{func() string { return L.ColUser }, 150},
	{func() string { return L.ColDisplay }, 220},
	{func() string { return L.ColAdded }, 215},
}

func (u *ui) memberCellText(row, col int) string {
	if row < 0 || row >= len(u.memList) {
		return ""
	}
	m := u.memList[row]
	switch col {
	case 0:
		return m.Role
	case 1:
		return m.Username
	case 2:
		return m.DisplayName
	case 3:
		return m.AddedAt
	}
	return ""
}

// membersPane bygger detalje-panelet under gridden.
func (u *ui) membersPane() fyne.CanvasObject {
	u.memTitle = widget.NewLabelWithStyle(L.MembersTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	u.memInfo = widget.NewLabel(L.SelectToSeeMembers)

	var gridBody fyne.CanvasObject
	u.memGrid, gridBody = newGridList(
		memberColumns,
		func() int { return len(u.memList) },
		u.memberCellText,
		nil, // ingen handling på valg — kun visning, men hele rækken markeres
		nil,
	)

	head := container.NewVBox(u.memTitle, u.memInfo)
	return container.NewBorder(head, nil, nil, nil, gridBody)
}

// clearMembers tømmer detalje-panelet (ingen database valgt).
func (u *ui) clearMembers() {
	u.memList = nil
	if u.memTitle != nil {
		u.memTitle.SetText(L.MembersTitle)
		u.memInfo.SetText(L.SelectToSeeMembers)
		u.memGrid.UnselectAll()
		u.memGrid.Refresh()
	}
}

// loadMembers henter medlemmerne for den valgte database og fylder panelet.
func (u *ui) loadMembers() {
	if u.selDB < 0 || u.selDB >= len(u.dbList) || u.memGrid == nil {
		return
	}
	db := u.dbList[u.selDB]
	u.memTitle.SetText(fmt.Sprintf(L.MembersOf, db.Name))
	u.memInfo.SetText(L.Working)
	u.memList = nil
	u.memGrid.UnselectAll()
	u.memGrid.Refresh()

	if !db.Bound {
		// `shares` slår op via lokal config — virker kun for bundne databaser.
		u.memInfo.SetText(L.MembersNeedBound)
		return
	}

	name := db.Name
	u.async(func() any {
		ctx, cancel := withTimeout(30 * time.Second)
		defer cancel()
		ms, r := u.c.shares(ctx, name)
		return memResult{ms: ms, r: r}
	}, func(v any) {
		res := v.(memResult)
		// Brugeren kan have valgt en anden række imens — ignorér forældet svar.
		if u.selDB < 0 || u.selDB >= len(u.dbList) || u.dbList[u.selDB].Name != name {
			return
		}
		if res.r.Err != nil {
			u.memList = nil
			u.memInfo.SetText(L.MembersUnavailable + " " + describeErr(res.r))
			u.memGrid.Refresh()
			return
		}
		u.memList = res.ms
		u.memInfo.SetText(fmt.Sprintf(L.MemberCount, len(res.ms)))
		u.memGrid.Refresh()
	})
}

// memResult bærer medlemslisten tilbage fra goroutinen.
type memResult struct {
	ms []member
	r  result
}

// syncSelected synkroniserer den valgte database i gridden.
func (u *ui) syncSelected() {
	if u.selDB < 0 || u.selDB >= len(u.dbList) {
		dialog.ShowInformation(L.Sync, L.SelectFirst, u.win)
		return
	}
	db := u.dbList[u.selDB]
	if !db.Bound {
		// Kun lokalt bundne databaser kan synkroniseres direkte.
		dialog.ShowInformation(L.Sync, L.OnServerOnly, u.win)
		return
	}
	u.promptSync(db.Name)
}

// forgetSelected fjerner den lokale binding for den valgte database via
// `forget <name>`. Det rører hverken .kdbx-filen eller databasen på serveren —
// kun koblingen i denne klients config.
func (u *ui) forgetSelected() {
	if u.selDB < 0 || u.selDB >= len(u.dbList) {
		dialog.ShowInformation(L.ForgetDatabase, L.SelectFirst, u.win)
		return
	}
	db := u.dbList[u.selDB]
	if !db.Bound {
		// En "kun på server"-database har ingen lokal binding at glemme.
		dialog.ShowInformation(L.ForgetDatabase, L.OnServerOnly, u.win)
		return
	}
	msg := fmt.Sprintf(L.ConfirmForget, db.Name)
	dialog.ShowConfirm(L.ForgetDatabase, msg, func(ok bool) {
		if !ok {
			return
		}
		name := db.Name
		u.async(func() any {
			ctx, cancel := withTimeout(30 * time.Second)
			defer cancel()
			return u.c.run(ctx, "", "forget", name)
		}, func(v any) {
			r := v.(result)
			u.log(r.Combined())
			if r.Err != nil {
				dialog.ShowError(errSimple(describeErr(r)), u.win)
				return
			}
			u.refreshDatabases()
		})
	}, u.win)
}

// deviceColumns er kolonnerne i enheds-gridden — samme som `devices`-kommandoen.
var deviceColumns = []gridColumn{
	{func() string { return L.ColStatus }, 110},
	{func() string { return L.ColName }, 160},
	{func() string { return L.ColID }, 290},
	{func() string { return L.ColEnrolled }, 215},
	{func() string { return L.ColLastSeen }, 215},
}

func (u *ui) deviceCellText(row, col int) string {
	if row < 0 || row >= len(u.devList) {
		return ""
	}
	d := u.devList[row]
	switch col {
	case 0:
		if d.Current {
			return L.ThisDevice
		}
		return ""
	case 1:
		return d.Name
	case 2:
		return d.ID
	case 3:
		return d.Enrolled
	case 4:
		return d.LastSeen
	}
	return ""
}

// devicesTab viser kontoens tilmeldte enheder som en grid (kontobredt — enheder
// hører til kontoen, ikke til en enkelt database).
func (u *ui) devicesTab() fyne.CanvasObject {
	u.devInfo = widget.NewLabel("")

	var gridBody fyne.CanvasObject
	u.devGrid, gridBody = newGridList(
		deviceColumns,
		func() int { return len(u.devList) },
		u.deviceCellText,
		func(row int) { u.selDev = row },
		func() { u.selDev = -1 },
	)

	add := widget.NewButton(L.AddDevice, func() { u.showAddDeviceDialog() })
	add.Importance = widget.HighImportance
	remove := widget.NewButton(L.RemoveDevice, func() { u.removeSelectedDevice() })
	refresh := widget.NewButton(L.Refresh, func() { u.refreshDevices() })

	toolbar := container.NewHBox(add, remove, layout.NewSpacer(), refresh)
	top := container.NewVBox(toolbar, u.devInfo)
	return container.NewBorder(top, nil, nil, nil, gridBody)
}

// showFormDialog viser en form-dialog der er markant bredere end Fynes default
// (som ellers krymper til indholdets bredde og bliver for smal til at indtaste
// stier, tokens og navne). Bredden er fast; højden skaleres med antal felter.
func (u *ui) showFormDialog(title, confirm string, items []*widget.FormItem, cb func(bool)) {
	d := dialog.NewForm(title, confirm, L.Cancel, items, cb, u.win)
	d.Resize(fyne.NewSize(620, float32(200+len(items)*48)))
	d.Show()
}

// showAddDeviceDialog (A) genererer en enrollment-token til en ny enhed via
// `admin user-enrollment <bruger> --admin-token <token>`. Enheder self-enroller,
// så GUI'en kan ikke tilmelde en fjern-maskine direkte — den udsteder en token
// som så bruges med `enroll` PÅ den nye enhed. Kræver en admin-token.
func (u *ui) showAddDeviceDialog() {
	username := widget.NewEntry()
	username.SetText(u.currentUser)
	username.SetPlaceHolder(L.Username)
	adminTok := widget.NewPasswordEntry()
	adminTok.SetPlaceHolder(L.AdminToken)

	items := []*widget.FormItem{
		widget.NewFormItem(L.Username, username),
		widget.NewFormItem(L.AdminToken, adminTok),
	}
	u.showFormDialog(L.AddDevice, L.AddDeviceCreate, items, func(ok bool) {
		if !ok {
			return
		}
		if username.Text == "" || adminTok.Text == "" {
			dialog.ShowError(errSimple(L.Username+" / "+L.AdminToken), u.win)
			return
		}
		u.async(func() any {
			ctx, cancel := withTimeout(30 * time.Second)
			defer cancel()
			return u.c.run(ctx, "", "admin", "user-enrollment", username.Text, "--admin-token", adminTok.Text)
		}, func(v any) {
			r := v.(result)
			u.log(r.Combined())
			if r.Err != nil {
				dialog.ShowError(errSimple(describeErr(r)), u.win)
				return
			}
			u.showEnrollTokenResult(r.Stdout)
		})
	})
}

// showEnrollTokenResult viser den genererede enrollment-token i et felt man kan
// kopiere fra.
func (u *ui) showEnrollTokenResult(output string) {
	out := widget.NewMultiLineEntry()
	out.SetText(output)
	out.Wrapping = fyne.TextWrapWord
	d := dialog.NewCustom(L.EnrollTokenCreated, L.Close, container.NewVScroll(out), u.win)
	d.Resize(fyne.NewSize(560, 320))
	d.Show()
}

// removeSelectedDevice (B) tilbagekalder den valgte enhed via `devices remove`.
func (u *ui) removeSelectedDevice() {
	if u.selDev < 0 || u.selDev >= len(u.devList) {
		dialog.ShowInformation(L.RemoveDevice, L.SelectDeviceFirst, u.win)
		return
	}
	d := u.devList[u.selDev]
	if d.Current {
		// Forhindr at man låser sig selv ude fra den enhed GUI'en kører på.
		dialog.ShowInformation(L.RemoveDevice, L.CannotRemoveCurrent, u.win)
		return
	}
	msg := fmt.Sprintf(L.ConfirmRemoveDevice, d.Name)
	dialog.ShowConfirm(L.RemoveDevice, msg, func(ok bool) {
		if !ok {
			return
		}
		id := d.ID
		u.async(func() any {
			ctx, cancel := withTimeout(30 * time.Second)
			defer cancel()
			return u.c.run(ctx, "", "devices", "remove", id)
		}, func(v any) {
			r := v.(result)
			u.log(r.Combined())
			if r.Err != nil {
				dialog.ShowError(errSimple(describeErr(r)), u.win)
				return
			}
			u.refreshDevices()
		})
	}, u.win)
}

// refreshDevices henter `devices` og opdaterer enheds-gridden.
func (u *ui) refreshDevices() {
	if u.devGrid == nil {
		return
	}
	u.devInfo.SetText(L.Working)
	u.async(func() any {
		ctx, cancel := withTimeout(30 * time.Second)
		defer cancel()
		ds, r := u.c.devices(ctx)
		return devResult{ds: ds, r: r}
	}, func(v any) {
		res := v.(devResult)
		u.selDev = -1
		u.devGrid.UnselectAll()
		if res.r.Err != nil {
			u.devList = nil
			u.log(res.r.Combined())
			u.devInfo.SetText(L.Error + ": " + describeErr(res.r))
			u.devGrid.Refresh()
			return
		}
		u.devList = res.ds
		u.devInfo.SetText(fmt.Sprintf(L.DevCount, len(res.ds)))
		u.log(fmt.Sprintf(L.DevCount, len(res.ds)))
		u.devGrid.Refresh()
	})
}

// devResult bærer enhedslisten tilbage fra goroutinen.
type devResult struct {
	ds []device
	r  result
}

// activityTab viser CLI-output-loggen. Log-feltet fylder hele fanen (et
// multiline-Entry har sin egen scrollbar), med en header og en Ryd-knap.
func (u *ui) activityTab() fyne.CanvasObject {
	clear := widget.NewButton(L.Clear, func() {
		u.activity.SetText("")
		u.activity.Refresh()
	})
	header := container.NewBorder(nil, nil, widget.NewLabel(L.ActivityHint), clear, nil)
	return container.NewBorder(header, nil, nil, nil, u.activity)
}

// settingsTab samler status, sprog, CLI-sti og nulstilling.
func (u *ui) settingsTab() fyne.CanvasObject {
	statusCard := widget.NewCard(L.StatusBox, "", u.statusLb)

	// VIGTIGT: sæt .Selected som felt FØR vi tildeler OnChanged. SetSelected()
	// ville udløse callbacken — og da callbacken genopbygger hele dashboardet,
	// gav den indledende værdisætning en uendelig løkke (dashboard byggede sig
	// selv igen og igen). Vi sætter derfor startværdien stille og reagerer kun
	// på reelle ændringer.
	langSel := widget.NewSelect([]string{"Dansk", "English"}, nil)
	if u.set.Language == string(langEN) {
		langSel.Selected = "English"
	} else {
		langSel.Selected = "Dansk"
	}
	langSel.OnChanged = func(choice string) {
		l := langDA
		if choice == "English" {
			l = langEN
		}
		if string(l) == u.set.Language {
			return // ingen reel ændring
		}
		setLang(l)
		u.set.Language = string(l)
		_ = saveSettings(u.set)
		u.showDashboard() // genopbyg med nye tekster
	}

	cliPath := widget.NewEntry()
	cliPath.SetText(u.c.path)
	cliBrowse := widget.NewButton(L.Browse, func() {
		dialog.ShowFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			defer rc.Close()
			p := uriToPath(rc.URI())
			cliPath.SetText(p)
			u.c.path = p
			u.set.CLIPath = p
			_ = saveSettings(u.set)
		}, u.win)
	})
	cliRow := container.NewBorder(nil, nil, nil, cliBrowse, cliPath)

	form := widget.NewForm(
		widget.NewFormItem(L.Language, langSel),
		widget.NewFormItem(L.CLIPathLabel, cliRow),
	)

	return container.NewVScroll(container.NewVBox(
		statusCard,
		widget.NewSeparator(),
		form,
	))
}

// refreshStatus henter `status` og opdaterer status-labelen.
func (u *ui) refreshStatus() {
	u.async(func() any {
		ctx, cancel := withTimeout(30 * time.Second)
		defer cancel()
		return u.c.status(ctx)
	}, func(v any) {
		s := v.(status)
		u.currentUser = usernameOnly(s.User)
		if u.statusLb == nil {
			return
		}
		if !s.Enrolled {
			u.statusLb.SetText(L.NotEnrolled)
			return
		}
		u.statusLb.SetText(
			"Server: " + s.Server + "\n" +
				"User: " + s.User + "\n" +
				"Device: " + s.Device + "\n" +
				"Last seen: " + s.LastSeen,
		)
	})
}

// refreshDatabases henter `databases` og opdaterer gridden.
func (u *ui) refreshDatabases() {
	if u.dbGrid == nil {
		return
	}
	u.dbInfo.SetText(L.Working)

	u.async(func() any {
		ctx, cancel := withTimeout(30 * time.Second)
		defer cancel()
		dbs, r := u.c.databases(ctx)
		debugf("databases: n=%d err=%v stderr=%q", len(dbs), r.Err, r.Stderr)
		return dbResult{dbs: dbs, r: r}
	}, func(v any) {
		res := v.(dbResult)
		debugf("databases done: n=%d", len(res.dbs))
		u.selDB = -1
		u.dbGrid.UnselectAll()
		u.clearMembers()
		if res.r.Err != nil {
			u.dbList = nil
			u.log(res.r.Combined())
			u.dbInfo.SetText(L.Error + ": " + describeErr(res.r))
			u.dbGrid.Refresh()
			return
		}
		u.dbList = res.dbs
		if len(res.dbs) == 0 {
			u.dbInfo.SetText(L.NoDatabases)
		} else {
			u.dbInfo.SetText(fmt.Sprintf(L.DBCount, len(res.dbs)))
		}
		u.log(fmt.Sprintf(L.DBCount, len(res.dbs)))
		u.dbGrid.Refresh()
	})
}

// dbResult bærer både listen og det rå resultat tilbage fra goroutinen.
type dbResult struct {
	dbs []database
	r   result
}

// promptSync beder om masterpassword og kører derefter sync for databasen.
func (u *ui) promptSync(name string) {
	pw := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem(L.MasterPwd, pw),
	}
	u.showFormDialog(L.MasterPwdFor+" "+name, L.Sync, items, func(ok bool) {
		if !ok || pw.Text == "" {
			return
		}
		u.runSync(name, pw.Text)
	})
}

// runSync udfører selve sync-kaldet og logger resultatet.
func (u *ui) runSync(name, masterPassword string) {
	u.log("sync " + name + " …")
	u.async(func() any {
		ctx, cancel := withTimeout(10 * time.Minute)
		defer cancel()
		return u.c.sync(ctx, name, masterPassword)
	}, func(v any) {
		r := v.(result)
		u.log(r.Combined())
		if r.Err != nil {
			dialog.ShowError(errSimple(describeErr(r)), u.win)
			return
		}
		u.refreshDatabases()
	})
}

// syncAll spørger om ét masterpassword og synkroniserer alle bundne databaser.
// Forudsætter at de deler password — ellers kan brugeren synkronisere én ad
// gangen via rækkens egen knap.
func (u *ui) syncAll() {
	pw := widget.NewPasswordEntry()
	items := []*widget.FormItem{widget.NewFormItem(L.MasterPwd, pw)}
	u.showFormDialog(L.SyncAll, L.Sync, items, func(ok bool) {
		if !ok || pw.Text == "" {
			return
		}
		u.async(func() any {
			ctx, cancel := withTimeout(30 * time.Second)
			defer cancel()
			dbs, _ := u.c.databases(ctx)
			return dbResult{dbs: dbs}
		}, func(v any) {
			for _, db := range v.(dbResult).dbs {
				if db.Bound {
					u.runSync(db.Name, pw.Text)
				}
			}
		})
	})
}

// showAddDBDialog registrerer en ny lokal database fra dashboardet.
func (u *ui) showAddDBDialog() {
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

	items := []*widget.FormItem{
		widget.NewFormItem(L.DBName, name),
		widget.NewFormItem(L.KdbxFile, pathRow),
	}
	u.showFormDialog(L.AddDatabase, L.CreateDBButton, items, func(ok bool) {
		if !ok || name.Text == "" || path.Text == "" {
			return
		}
		u.async(func() any {
			ctx, cancel := withTimeout(60 * time.Second)
			defer cancel()
			return u.c.initDB(ctx, name.Text, path.Text)
		}, func(v any) {
			r := v.(result)
			u.log(r.Combined())
			if r.Err != nil {
				dialog.ShowError(errSimple(describeErr(r)), u.win)
				return
			}
			u.refreshDatabases()
		})
	})
}
