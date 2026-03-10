package model

import "syscall/js"

type Theme struct {
	Path string
}

func NewTheme(path string) *Theme {
	return &Theme{
		Path: path,
	}
}

func (t *Theme) AddStylesheet(id, href string) {
	doc := js.Global().Get("document")

	// avoid duplicates
	existing := doc.Call("getElementById", id)
	if !existing.IsNull() {
		return
	}

	link := doc.Call("createElement", "link")
	link.Set("id", id)
	link.Set("rel", "stylesheet")
	link.Set("href", href)
	doc.Get("head").Call("appendChild", link)
}

func (t *Theme) RemoveStylesheet(id string) {
	doc := js.Global().Get("document")
	link := doc.Call("getElementById", id)
	if !link.IsNull() {
		link.Get("parentNode").Call("removeChild", link)
	}
}

func (t *Theme) SetTheme(path string) {
	t.Path = path
	t.RemoveStylesheet("theme-stylesheet")
	t.AddStylesheet("theme-stylesheet", path+".css")
}
