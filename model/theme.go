package model

import (
	"fmt"
	"syscall/js"
)

var (
	// i sure do love these global variables
	ThemChan chan string = make(chan string, 3)
	// GreenFlag bool = false
	// there are also technically race conditions with all of these but like pfffft

	CurrentPalette *Palette = &CatppuccinPalette

	Palettes = map[string]Palette{
		"catppuccin": CatppuccinPalette,
		"dark":       DarkPalette,
	}

	CatppuccinPalette = Palette{
		Base:   "#1e1e2e",
		Mantle: "#181825",
		Crust:  "#11111b",

		Surface0: "#313244",
		Surface1: "#45475a",
		Surface2: "#585b70",

		Overlay0: "#6c7086",
		Overlay1: "#7f849c",
		Overlay2: "#9399b2",

		Text:     "#cdd6f4",
		Subtext0: "#a6adc8",
		Subtext1: "#bac2de",

		Rosewater: "#f5e0dc",
		Flamingo:  "#f2cdcd",
		Pink:      "#f5c2e7",
		Mauve:     "#cba6f7",
		Red:       "#f38ba8",
		Maroon:    "#eba0ac",
		Peach:     "#fab387",
		Yellow:    "#f9e2af",
		Green:     "#a6e3a1",
		Teal:      "#94e2d5",
		Sky:       "#89dceb",
		Sapphire:  "#74c7ec",
		Blue:      "#89b4fa",
		Lavender:  "#b4befe",
	}

	DarkPalette = Palette{
		Base:   "#1e1e1e",
		Mantle: "#181818",
		Crust:  "#121212",

		Surface0: "#2a2a2a",
		Surface1: "#333333",
		Surface2: "#3d3d3d",

		Overlay0: "#555555",
		Overlay1: "#666666",
		Overlay2: "#777777",

		Text:     "#e0e0e0",
		Subtext0: "#9e9e9e",
		Subtext1: "#bdbdbd",

		Rosewater: "#ffccbc",
		Flamingo:  "#ffab91",
		Pink:      "#f48fb1",
		Mauve:     "#ce93d8",
		Red:       "#e57373",
		Maroon:    "#ef9a9a",
		Peach:     "#ffb74d",
		Yellow:    "#fff176",
		Green:     "#81c784",
		Teal:      "#4db6ac",
		Sky:       "#4fc3f7",
		Sapphire:  "#29b6f6",
		Blue:      "#64b5f6",
		Lavender:  "#9fa8da",
	}
)

type Palette struct {
	Base   string
	Mantle string
	Crust  string

	// surfaces
	Surface0 string
	Surface1 string
	Surface2 string

	// overlays
	Overlay0 string
	Overlay1 string
	Overlay2 string

	// text
	Text     string
	Subtext0 string
	Subtext1 string

	// accents
	Rosewater string
	Flamingo  string
	Pink      string
	Mauve     string
	Red       string
	Maroon    string
	Peach     string
	Yellow    string
	Green     string
	Teal      string
	Sky       string
	Sapphire  string
	Blue      string
	Lavender  string
}

// wow i bet it sure did take a long time to make these structs ...

func (p Palette) ToCSS() string {
	return fmt.Sprintf(`
		:root {
			--base:      %s;
			--mantle:    %s;
			--crust:     %s;

			--surface0:  %s;
			--surface1:  %s;
			--surface2:  %s;

			--overlay0:  %s;
			--overlay1:  %s;
			--overlay2:  %s;

			--text:      %s;
			--subtext0:  %s;
			--subtext1:  %s;

			--rosewater: %s;
			--flamingo:  %s;
			--pink:      %s;
			--mauve:     %s;
			--red:       %s;
			--maroon:    %s;
			--peach:     %s;
			--yellow:    %s;
			--green:     %s;
			--teal:      %s;
			--sky:       %s;
			--sapphire:  %s;
			--blue:      %s;
			--lavender:  %s;
		}
	`,
		p.Base, p.Mantle, p.Crust,
		p.Surface0, p.Surface1, p.Surface2,
		p.Overlay0, p.Overlay1, p.Overlay2,
		p.Text, p.Subtext0, p.Subtext1,
		p.Rosewater, p.Flamingo, p.Pink, p.Mauve,
		p.Red, p.Maroon, p.Peach, p.Yellow,
		p.Green, p.Teal, p.Sky, p.Sapphire,
		p.Blue, p.Lavender,
	)
}

type Theme struct {
	Name string
}

func NewTheme(path string) *Theme {
	return &Theme{
		Name: path,
	}
}

func (p *Palette) Apply() {
	doc := js.Global().Get("document")

	// avoid duplicates
	// existing := doc.Call("getElementById", "theme")
	// if !existing.IsNull() {
	// 	fmt.Println("this is dumb")
	// 	return
	// }

	style := doc.Call("createElement", "style")
	style.Set("id", "theme")
	style.Set("innerHTML", p.ToCSS())
	doc.Get("head").Call("appendChild", style)
}

func (t *Theme) RemoveStylesheet(id string) {
	doc := js.Global().Get("document")
	link := doc.Call("getElementById", id)
	if !link.IsNull() {
		link.Get("parentNode").Call("removeChild", link)
	}
}

func (t *Theme) SetTheme(name string) {

	if name == t.Name {
		fmt.Println("same theme; not changing...")
		return
	}

	if palette, ok := Palettes[name]; ok {
		t.Name = name
		CurrentPalette = &palette

		palette.Apply()
	} else {
		fmt.Println("could not find theme name")
	}

}
