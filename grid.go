// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Fyne's widget.Table markerer én CELLE ad gangen — der findes ingen native
// hel-række-markering. Vores grids bygges derfor på widget.List, som markerer
// hele rækken. Hver liste-række er en vandret stribe celler med faste
// kolonnebredder (columnsLayout), og en matchende header ligger ovenover.

// gridColumn beskriver én kolonne: overskrift (funktion så sprog-skift virker)
// og fast pixelbredde.
type gridColumn struct {
	header func() string
	width  float32
}

// columnsLayout lægger børn vandret ud med en fast bredde per kolonne, så
// header og alle rækker flugter. Højden er den højeste celles højde.
type columnsLayout struct {
	widths []float32
}

func (c *columnsLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	var w, h float32
	for i, o := range objs {
		if i < len(c.widths) {
			w += c.widths[i]
		}
		if ms := o.MinSize(); ms.Height > h {
			h = ms.Height
		}
	}
	return fyne.NewSize(w, h)
}

func (c *columnsLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	var x float32
	for i, o := range objs {
		w := float32(0)
		if i < len(c.widths) {
			w = c.widths[i]
		}
		o.Resize(fyne.NewSize(w, size.Height))
		o.Move(fyne.NewPos(x, 0))
		x += w
	}
}

// newGridList bygger en række-valgbar, kolonne-justeret grid med header.
// Returnerer både listen (til Refresh/UnselectAll) og det færdige widget (med
// header + vandret scroll til brede kolonner).
func newGridList(
	cols []gridColumn,
	length func() int,
	cellText func(row, col int) string,
	onSelect func(row int),
	onUnselect func(),
) (*widget.List, fyne.CanvasObject) {
	widths := make([]float32, len(cols))
	for i, c := range cols {
		widths[i] = c.width
	}
	lay := &columnsLayout{widths: widths}

	hdr := make([]fyne.CanvasObject, len(cols))
	for i, c := range cols {
		l := widget.NewLabel(c.header())
		l.TextStyle.Bold = true
		l.Truncation = fyne.TextTruncateEllipsis
		hdr[i] = l
	}
	header := container.New(lay, hdr...)

	list := widget.NewList(
		length,
		func() fyne.CanvasObject {
			objs := make([]fyne.CanvasObject, len(cols))
			for i := range cols {
				l := widget.NewLabel("")
				l.Truncation = fyne.TextTruncateEllipsis
				objs[i] = l
			}
			return container.New(lay, objs...)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			for i := range cols {
				row.Objects[i].(*widget.Label).SetText(cellText(id, i))
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		if onSelect != nil {
			onSelect(id)
		}
	}
	list.OnUnselected = func(widget.ListItemID) {
		if onUnselect != nil {
			onUnselect()
		}
	}

	// Header øverst, listen fylder resten; vandret scroll når kolonnerne er
	// bredere end vinduet.
	body := container.NewBorder(header, nil, nil, nil, list)
	return list, container.NewHScroll(body)
}
