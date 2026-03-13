package components

import (
	"strconv"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"

	"github.com/dq1Mango/graph-theory/actions"
)

type EditableButton struct {
	vecty.Core
	label   string
	editing bool
	draft   string

	// Submit chan<- string
	onSubmit func(string)
}

func NewEditableButton(label string, onSubmit func(string)) *EditableButton {
	return &EditableButton{
		label: label,
		draft: "",

		onSubmit: onSubmit,
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
		e.onSubmit(e.draft)
	}
	e.editing = false
	vecty.Rerender(e)
}

func (e *EditableButton) cancel() {
	e.draft = ""
	e.editing = false
	vecty.Rerender(e)
}

func NewGraphFromPrufer(action chan any) *EditableButton {

	// pruferCodeGeneration := make(chan string, 2)

	return NewEditableButton(
		"construct graph from prufer code",

		func(sequence string) {

			var paresed []uint

			for segment := range strings.SplitSeq(sequence, ",") {
				segment = strings.TrimSpace(segment)

				num, err := strconv.Atoi(segment)

				if err != nil {
					InfoChan <- "Cannot Parse Prufer Code"

					// 'goto considered harmful' -Dijkstra 1968
					// goto start
				}

				paresed = append(paresed, uint(num))
			}

			action <- &actions.GraphFromPrufer{Prufer: paresed}
		},
	)

}

// dang i had this rly fun way to do it with gotos, but i think its just worse
// func (p *GraphFromPrufer) Mount() {
//
// 	go func() {
// 		for {
// 		start:
//
// 			sequence := <-p.Submit
// 			var paresed []uint
//
// 			for segment := range strings.SplitSeq(sequence, ",") {
// 				segment = strings.TrimSpace(segment)
//
// 				num, err := strconv.Atoi(segment)
//
// 				if err != nil {
// 					InfoChan <- "Cannot Parse Prufer Code"
//
// 					// 'goto considered harmful' -Dijkstra 1968
// 					goto start
// 				}
//
// 				paresed = append(paresed, uint(num))
// 			}
//
// 			p.actions <- &actions.GraphFromPrufer{Prufer: paresed}
// 		}
// 	}()
// }
