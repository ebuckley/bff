package bff

import (
	"context"
	"errors"
)

var ErrActionAlreadyExists = errors.New("action already exists")
var ErrActionNotFound = errors.New("action not found")

type HandlerFunc func(ctx context.Context, io *Io) error

type Action struct {
	handler     HandlerFunc
	display     chan string
	input       chan any
	Slug        string `json:"slug,omitempty"`
	Name        string `json:"name" json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ActionOption func(*Action)

func NewAction(name string, handler HandlerFunc, opts ...ActionOption) *Action {
	action := &Action{
		Name:    name,
		Slug:    name,
		handler: handler,
	}
	for _, opt := range opts {
		opt(action)
	}
	return action
}

func WithSlug(slug string) ActionOption {
	return func(a *Action) {
		a.Slug = slug
	}
}
func WithDescription(description string) ActionOption {
	return func(a *Action) {
		a.Description = description
	}
}
