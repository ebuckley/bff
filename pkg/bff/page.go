package bff

import "errors"

var ErrPageAlreadyExists = errors.New("page already exists")

type Page struct {
	actions []*Action
	name    string
}

func NewPage(name string) *Page {
	return &Page{name: name}
}

func (p *Page) AddAction(action *Action) {
	p.actions = append(p.actions, action)
}
