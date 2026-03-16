package components

import (
	"github.com/dq1Mango/graph-theory/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type ThemeChanger struct {
	vecty.Core
	dropdown *Dropdown
}

func NewThemeChanger() *ThemeChanger {

	themeChange := func(newTheme string) {
		model.ThemChan <- newTheme
	}
	// i rly like the variadic functions
	themeDropDown := NewDropdown("theme", themeChange, "dark", "catppuccin")

	return &ThemeChanger{dropdown: themeDropDown}
}

func (t *ThemeChanger) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("theme-changer"),
		),
		t.dropdown,
	)
}
