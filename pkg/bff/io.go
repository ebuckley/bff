package bff

import (
	"fmt"
)

type Io struct {
	stack   []Executable
	Display Display
	Input   Input
	input   <-chan Message
	output  chan<- Message
}

func NewIo(input <-chan Message, output chan<- Message) *Io {
	io := &Io{
		stack:  make([]Executable, 0),
		input:  input,
		output: output,
	}
	display := Display{io}
	i := Input{io}
	io.Display = display
	io.Input = i
	return io
}

// Display represents the display device, call methods to add display content to the stack
type Display struct {
	io *Io
}

// Input represents the input device, call methods to add input requests to the stack
type Input struct {
	io *Io
}
type Executable interface {
	Execute(input <-chan Message, output chan<- Message) (any, error)
}

// need to support
// DONE:
// - Group: Combines multiple I/O method calls into a single form.
// - input.Text requests a string value
// - input.boolean requests a boolean value
// TODO:
// - input.number requests a number value
// - input.email requests an email value
// - input.slider requests a number value within a range
// - input.date requests a date value
// - input.richText requests a rich text value
// - input.url requests a URL value
// - input.time requests a time value
// - input.file requests a file value
// - input.confirm requests confirmation of an action using a full screen dialog box
// - input.confirmIdentity (multi factor with the users email)
// - input.search search for arbitrary results using a search box
// - input.selectTable requests a selection from a table of options
// - input.selectSingle Prompts the app user to select a single value from a set of provided values.

// InputBase defines everything that all inputs have in common
type InputBase struct {
	Label       string `json:"label,omitempty"`
	HelpText    string `json:"helpText,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// TextInput is a text box input
type TextInput struct {
	InputBase
	MinLength int `json:"minLength,omitempty"`
	MaxLength int `json:"maxLength,omitempty"`
}

func (h *TextInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "textInput", Data: h}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

type BooleanInput struct {
	InputBase
}

func (h *BooleanInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "booleanInput", Data: h}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

type NumberInput struct {
	InputBase
	Min float64
	Max float64
}

func (h *NumberInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "numberInput", Data: h}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

//---------------
//  Display Types
//---------------
// DONE:
// - display.heading Displays a heading to the action user.
// - display.markdown Displays rendered markdown to the action user. display.metadata

// TODO:
// - display.html Displays rendered HTML to the action user.
// - display.code Displays a block of code to the action user.
// - display.grid  Displays data in a grid layout.
// - display.image Displays an image to the action user. One of url or buffer must be provided.
// - display.link Displays a button-styled action link to the action user. Can link to an external URL or to another action.
// - display.metadata Displays a series of label/value pairs in a variety of layout options.
// - display.object Displays an object of nested data to the action user.
// - display.table Displays tabular data.
// - display.video Displays a video to the action user. One of url or buffer must be provided.

type CodeDisplay struct {
	Code     string `json:"code,omitempty"`
	Language string `json:"language,omitempty"`
}

func (c CodeDisplay) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "code", Data: c}
	return nil, nil
}

type HeadingDisplay struct {
	Text  string `json:"text,omitempty"`
	Level int    `json:"level,omitempty"`
}

func (h HeadingDisplay) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "display", Data: h}
	return nil, nil
}

type MarkdownDisplay struct {
	Content string `json:"content"`
}

func (m MarkdownDisplay) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "markdown", Data: m}
	return nil, nil
}

// Implement other input and display types similarly...

// Group combines multiple I/O method calls
type Group struct {
	Elements []Executable
}

func (g Group) Execute(input <-chan Message, output chan<- Message) error {
	output <- Message{Type: "group", Data: g}
	return nil
}

// AddToStack adds the element to the stack and executes it -- returning the result of the execution
func (io *Io) AddToStack(element Executable) (any, error) {
	io.stack = append(io.stack, element)
	return element.Execute(io.input, io.output)
}

func (d *Display) Heading(text string, level int) {
	d.io.AddToStack(HeadingDisplay{Text: text, Level: level})
}

func (d *Display) Code(code string, language string) {
	d.io.AddToStack(CodeDisplay{Code: code, Language: language})
}

func (d *Display) Markdown(content string) {
	d.io.AddToStack(MarkdownDisplay{Content: content})
}

func (i *Input) Text(label string, options ...func(*TextInput)) (string, error) {
	input := &TextInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", v)
	}
	return s, nil
}

func (i *Input) Boolean(label string, options ...func(*BooleanInput)) (bool, error) {
	input := &BooleanInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("expected boolean, got %T", v)
	}
	return b, nil
}

// Implement other input methods similarly...

// WithHelpText is an option function to set the help text of an input
func WithHelpText(text string) func(*InputBase) {
	return func(i *InputBase) {
		i.HelpText = text
	}
}

// WithPlaceholder is an option function to set the placeholder of an input
func WithPlaceholder(placeholder string) func(*InputBase) {
	return func(i *InputBase) {
		i.Placeholder = placeholder
	}
}

// WithRequired is an option function to set the required status of an input
func WithRequired(required bool) func(*InputBase) {
	return func(i *InputBase) {
		i.Required = required
	}
}

// Add more option functions as needed...
