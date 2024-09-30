package bff

import (
	"context"
	"sync"
)

// BFF represents the Backend for Frontend, which manages actions and pages
type BFF struct {
	actions     map[string]*Action
	pages       map[string]*Page
	mu          sync.RWMutex
	io          *Io
	environment string
}

// New creates a new BFF instance
func New(environment string) *BFF {
	return &BFF{
		actions:     make(map[string]*Action),
		pages:       make(map[string]*Page),
		io:          &Io{},
		environment: environment,
	}
}

// RegisterAction adds a new action to the BFF
func (b *BFF) RegisterAction(name string, handler HandlerFunc, opts ...ActionOption) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.actions[name]; exists {
		return ErrActionAlreadyExists
	}

	a := NewAction(name, handler, opts...)
	b.actions[name] = a
	return nil
}

// RegisterPage adds a new page to the BFF
func (b *BFF) RegisterPage(name string, actions []string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.pages[name]; exists {
		return ErrPageAlreadyExists
	}

	p := NewPage(name)
	for _, actionName := range actions {
		if action, exists := b.actions[actionName]; exists {
			p.AddAction(action)
		} else {
			return ErrActionNotFound
		}
	}

	b.pages[name] = p
	return nil
}

// ExecuteAction runs the specified action
func (b *BFF) ExecuteAction(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	b.mu.RLock()
	action, exists := b.actions[name]
	b.mu.RUnlock()

	if !exists {
		return nil, ErrActionNotFound
	}

	return action.Execute(ctx, b.io, params)
}

// GetPages returns all registered pages
func (b *BFF) GetPages() []*Page {
	b.mu.RLock()
	defer b.mu.RUnlock()

	pages := make([]*Page, 0, len(b.pages))
	for _, p := range b.pages {
		pages = append(pages, p)
	}
	return pages
}

// GetEnvironment returns the current environment
func (b *BFF) GetEnvironment() string {
	return b.environment
}

func (b *BFF) GetActions() []*Action {
	b.mu.RLock()
	defer b.mu.RUnlock()

	actions := make([]*Action, 0, len(b.actions))
	for _, a := range b.actions {
		actions = append(actions, a)
	}
	return actions
}
