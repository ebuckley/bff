package bff

import (
	"context"
	"log/slog"
	"sync"
)

// Message represents a message with the backend
type Message struct {
	Type string `json:"type,omitempty"`
	Data any    `json:"data,omitempty"`
}

// BFF represents the Backend for Frontend, which manages actions and pages
type BFF struct {
	actions     map[string]*Action
	pages       map[string]*Page
	mu          sync.RWMutex
	environment string
}

// New creates a new BFF instance
func New(environment string) *BFF {
	return &BFF{
		actions:     make(map[string]*Action),
		pages:       make(map[string]*Page),
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
func (b *BFF) ExecuteAction(ctx context.Context, name string, input <-chan Message, output chan<- Message) error {
	b.mu.RLock()
	action, exists := b.actions[name]
	b.mu.RUnlock()
	if !exists {
		return ErrActionNotFound
	}
	// make a nice little IO context we can give to the action to handle
	io := NewIo(input, output)

	return action.handler(ctx, io)
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

func (b *BFF) Loop(ctx context.Context, input <-chan Message, output chan<- Message) {
	// the application loop
	for {
		select {
		case <-ctx.Done():
			slog.Debug("exiting bff loop with connection")
			return
		case v := <-input:
			if v.Type == "start" {
				// pass the input/output chanel to execute action
				name, ok := v.Data.(string)
				if !ok {
					output <- Message{Type: "error", Data: "expected string"}
					continue
				}
				err := b.ExecuteAction(ctx, name, input, output)
				if err != nil {
					output <- Message{Type: "error", Data: err.Error()}
					slog.Error("failed to execute action: ", "err", err)
					return
				}
				return
			}
		}
	}
}
