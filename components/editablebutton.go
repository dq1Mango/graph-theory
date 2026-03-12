package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type EditableButton struct {
	vecty.Core
	label   string
	editing bool
	draft   string

	Submit chan<- string
}

func NewEditableButton(label string, submit chan<- string) *EditableButton {
	return &EditableButton{
		label: label,
		draft: "",

		Submit: submit,
	}
}

func (e *EditableButton) Render() vecty.ComponentOrHTML {
	if e.editing {
		return e.renderEditor()
	}
	return e.renderButton()
}

func (e *EditableButton) renderButton() vecty.ComponentOrHTML {
	return elem.Button(
		vecty.Markup(
			vecty.Class("editable-button"),
			event.Click(func(ev *vecty.Event) {
				e.editing = true
				vecty.Rerender(e)
			}),
		),
		vecty.Text(e.label),
	)
}

func (e *EditableButton) renderEditor() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Input(
			vecty.Markup(
				vecty.Class("editable-button-input"),
				vecty.Attribute("type", "text"),
				vecty.Attribute("value", e.draft),
				event.Input(func(ev *vecty.Event) {
					e.draft = ev.Value.Get("target").Get("value").String()
				}),
				event.KeyDown(func(ev *vecty.Event) {
					if ev.Value.Get("key").String() == "Enter" {
						e.commit()
					} else if ev.Value.Get("key").String() == "Escape" {
						e.cancel()
					}
				}),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("editable-button-editor"),
			),
			elem.Button(
				vecty.Markup(
					vecty.Class("editable-button-submit"),
					event.Click(func(ev *vecty.Event) {
						e.commit()
					}),
				),
				vecty.Text("✓"),
			),
			elem.Button(
				vecty.Markup(
					vecty.Class("editable-button-cancel"),
					event.Click(func(ev *vecty.Event) {
						e.cancel()
					}),
				),
				vecty.Text("✕"),
			),
		))
}

func (e *EditableButton) commit() {
	if e.draft != "" {
		e.Submit <- e.draft
	}
	e.editing = false
	vecty.Rerender(e)
}

func (e *EditableButton) cancel() {
	e.draft = ""
	e.editing = false
	vecty.Rerender(e)
}
