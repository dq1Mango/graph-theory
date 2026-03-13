package components

import (
	"syscall/js"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
)

const DISPLAY_TIME = 5
const ID = "popup"

// u can only make one of these cuze of this little guy
var InfoChan = make(chan string, 3)

type Popup struct {
	vecty.Core
	element js.Value

	Message       string
	AnimationTime uint

	animating bool
	interrupt chan struct{}
}

// func (p *Popup) Clear() {
// 	p.animating = false
// 	p.ctx.Set("textContent", "")
// }

func (p *Popup) FadeOut() {

	p.animating = true
	for p.animating {

		vecty.Rerender(p)
		select {

		case <-time.After(5 * time.Second):
			p.animating = false
			p.Message = ""
			vecty.Rerender(p)

		case <-p.interrupt:

			p.element.Get("classList").Call("remove", "fadeout")
			p.element.Get("offsetWidth") // force reflow
			p.element.Get("classList").Call("add", "fadeout")

		}
	}
}

func (p *Popup) ListenForUpdates() {
	for {
		newMessage := <-InfoChan

		p.Message = newMessage

		if p.animating {
			p.interrupt <- struct{}{}
		} else {
			go p.FadeOut()
		}
	}
}

func (p *Popup) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Paragraph(
			vecty.Text(p.Message),
			vecty.Markup(
				prop.ID(ID),
				vecty.ClassMap{
					"fadeout": p.animating,
					"popup":   true,
				}),
		),
	)
}

func (p *Popup) Mount() {

	p.element = js.Global().Get("document").Call("getElementById", ID)

	go p.ListenForUpdates()

	// if p.ctx.IsUndefined() || p.ctx.IsNull() || p.ctx.IsNull() {
	// 	fmt.Println("heres the problem")
	// }
	// fmt.Println("really?")
}

func NewPopup() Popup {
	popup := Popup{
		Message:       "",
		AnimationTime: 5,

		interrupt: make(chan struct{}),
	}

	return popup
}
