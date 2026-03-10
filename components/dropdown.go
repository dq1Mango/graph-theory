package components

import (
	"github.com/dq1Mango/graph-theory/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	// "github.com/hexops/vecty/prop"
)

type Dropdown struct {
	vecty.Core

	title    string
	options  []string
	selected string
	open     bool
}

func NewDropdown(title string, options ...string) *Dropdown {
	return &Dropdown{
		title:    title,
		options:  options,
		selected: options[0],
	}
}

func (d *Dropdown) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("dropdown")),
		d.renderTrigger(),
		d.renderMenu(),
	)
}

func (d *Dropdown) renderTrigger() vecty.ComponentOrHTML {
	return elem.Button(
		vecty.Markup(
			vecty.Class("dropdown-button"),
			event.Click(func(e *vecty.Event) {
				d.open = !d.open
				vecty.Rerender(d)
			}),
		),
		vecty.Text(d.title+" ▾"),
	)
}

func (d *Dropdown) renderMenu() vecty.ComponentOrHTML {

	items := make(vecty.List, len(d.options))
	for i, opt := range d.options {
		opt := opt // capture loop variable
		items[i] = elem.Button(
			vecty.Markup(
				vecty.Class("dropdown-item"),
				vecty.ClassMap{
					"dropdown-item-selected": opt == d.selected,
				},
				event.Click(func(e *vecty.Event) {
					d.selected = opt
					d.open = false
					vecty.Rerender(d)

					model.ThemChan <- d.selected
				}),
			),
			vecty.Text(opt),
		)
	}

	return elem.Div(
		vecty.Markup(vecty.Class("dropdown-content")), //, prop.Disabled(!d.open)),
		items,
	)
}
