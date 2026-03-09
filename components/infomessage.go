package components

import (
	// "fmt"
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
	ctx js.Value

	Message  string
	TimeLeft uint

	animating bool
}

// func (p *Popup) Clear() {
// 	p.animating = false
// 	p.ctx.Set("textContent", "")
// }

func (p *Popup) FadeOut() {
	p.animating = true
	vecty.Rerender(p)

	time.Sleep(4 * time.Second)

	p.animating = false
	p.Message = ""
	vecty.Rerender(p)

	// p.Clear()
}

func (p *Popup) ListenForUpdates() {
	for {
		newMessage := <-InfoChan

		p.Message = newMessage
		p.TimeLeft = DISPLAY_TIME

		// p.ctx.Set("textContent", p.Message)

		go p.FadeOut()
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

	p.ctx = js.Global().Get("document").Call("getElementById", ID)

	go p.ListenForUpdates()

	// if p.ctx.IsUndefined() || p.ctx.IsNull() || p.ctx.IsNull() {
	// 	fmt.Println("heres the problem")
	// }
	// fmt.Println("really?")
}

func NewPopup() Popup {
	popup := Popup{
		Message:  "",
		TimeLeft: 0,
	}

	return popup
}
