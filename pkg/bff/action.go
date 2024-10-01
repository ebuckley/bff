package bff

import (
	"context"
	"errors"
	"fmt"
)

var ErrActionAlreadyExists = errors.New("action already exists")
var ErrActionNotFound = errors.New("action not found")

type HandlerFunc func(ctx context.Context, io *Io) error

type Action struct {
	handler     HandlerFunc
	display     chan string
	input       chan any
	Name        string `json:"name" json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ActionOption func(*Action)

func NewAction(name string, handler HandlerFunc, opts ...ActionOption) *Action {
	action := &Action{
		Name:    name,
		handler: handler,
	}
	for _, opt := range opts {
		opt(action)
	}
	return action
}

func (a *Action) Render() string {
	return fmt.Sprintf(`{"Name": "%s", "Description": "%s"}`, a.Name, a.Description)
}
