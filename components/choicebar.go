package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type ChoiceBar struct {
	vecty.Core
	Choices  []string
	Selected uint

	Clicks  chan uint
	Updates chan<- string
}

func NewChoiceBar(updates chan<- string, choices ...string) *ChoiceBar {
	return &ChoiceBar{
		Choices:  choices,
		Selected: 0,

		Clicks:  make(chan uint, 2),
		Updates: updates,
	}
}

func (c *ChoiceBar) Render() vecty.ComponentOrHTML {
	elements := make([]vecty.MarkupOrChild, len(c.Choices)+1)

	for i, choice := range c.Choices {
		i := uint(i)

		elements[i] = elem.Button(
			vecty.Markup(
				vecty.Class("choice"),

				vecty.ClassMap{
					"selected-choice": c.Selected == i,
				},

				event.Click(func(e *vecty.Event) {
					c.Clicks <- i
				}),
			),

			vecty.Text(choice),
		)
	}

	elements = append(elements,
		vecty.Markup(
			vecty.Class("choice-bar"),
		),
	)

	return elem.Div(
		elements...,
	)
}

func (c *ChoiceBar) Mount() {
	go func() {
		for {
			newChoice := <-c.Clicks

			if newChoice != c.Selected {
				c.Selected = newChoice

				c.Updates <- c.Choices[c.Selected]
				vecty.Rerender(c)
			}
		}
	}()
}
