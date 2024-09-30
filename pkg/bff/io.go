package bff

import (
	"fmt"
	"strings"
)

type Io struct {
	stack   []any
	Display Display
	Input   Input
}

// Display represents the display device, call methods to add display content to the stack
type Display struct {
	io *Io
}

// Input represents the input device, call methods to add input requests to the stack
type Input struct {
	io *Io
}

type Renderable interface {
	// Render returns the JSON representation state of the element, for showing something on the screen
	Render() string
}
type Inputable interface {
	// Input will wait for the input and return the value.. it means network round trips etc etc
	Input(any) (any, error)
}

// need to support
// - Group: Combines multiple I/O method calls into a single form.
// - input.Text requests a string value
// - input.boolean requests a boolean value
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
// - display.code Displays a block of code to the action user.
// - display.grid  Displays data in a grid layout.
// - display.heading Displays a heading to the action user.
// - display.html Displays rendered HTML to the action user.
// - display.image Displays an image to the action user. One of url or buffer must be provided.
// - display.link Displays a button-styled action link to the action user. Can link to an external URL or to another action.
// - display.markdown Displays rendered markdown to the action user. display.metadata
// - display.metadata Displays a series of label/value pairs in a variety of layout options.
// - display.object Displays an object of nested data to the action user.
// - display.table Displays tabular data.
// - display.video Displays a video to the action user. One of url or buffer must be provided.

// InputBase defines everything that all inputs have in common
type InputBase struct {
	Label       string
	HelpText    string
	Placeholder string
	Required    bool
}

// TextInput is a text box input
type TextInput struct {
	InputBase
	MinLength int
	MaxLength int
}

func (t TextInput) Render() string {
	return fmt.Sprintf(`{"type":"text","label":"%s","helpText":"%s","placeholder":"%s","required":%t,"minLength":%d,"maxLength":%d}`,
		t.Label, t.HelpText, t.Placeholder, t.Required, t.MinLength, t.MaxLength)
}

func (t TextInput) Input(value any) (any, error) {
	str, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", value)
	}
	if len(str) < t.MinLength || (t.MaxLength > 0 && len(str) > t.MaxLength) {
		return nil, fmt.Errorf("input length must be between %d and %d", t.MinLength, t.MaxLength)
	}
	return str, nil
}

// BooleanInput
type BooleanInput struct {
	InputBase
}

func (b BooleanInput) Render() string {
	return fmt.Sprintf(`{"type":"boolean","label":"%s","helpText":"%s","required":%t}`,
		b.Label, b.HelpText, b.Required)
}

func (b BooleanInput) Input(value any) (any, error) {
	boolVal, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("expected boolean, got %T", value)
	}
	return boolVal, nil
}

// NumberInput
type NumberInput struct {
	InputBase
	Min float64
	Max float64
}

func (n NumberInput) Render() string {
	return fmt.Sprintf(`{"type":"number","label":"%s","helpText":"%s","placeholder":"%s","required":%t,"min":%f,"max":%f}`,
		n.Label, n.HelpText, n.Placeholder, n.Required, n.Min, n.Max)
}

func (n NumberInput) Input(value any) (any, error) {
	num, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", value)
	}
	if num < n.Min || (n.Max > 0 && num > n.Max) {
		return nil, fmt.Errorf("input must be between %f and %f", n.Min, n.Max)
	}
	return num, nil
}

// Display types
type CodeDisplay struct {
	Code     string
	Language string
}

func (c CodeDisplay) Render() string {

	return fmt.Sprintf(`{"type":"code","code":%s,"language":"%s"}`,
		c.Code, c.Language)
}

type HeadingDisplay struct {
	Text  string
	Level int
}

func (h HeadingDisplay) Render() string {
	return fmt.Sprintf(`{"type":"heading","text":"%s","level":%d}`, h.Text, h.Level)
}

type MarkdownDisplay struct {
	Content string
}

func (m MarkdownDisplay) Render() string {
	return fmt.Sprintf(`{"type":"markdown","content":%s}`, m.Content)
}

// Implement other input and display types similarly...

// Group combines multiple I/O method calls
type Group struct {
	Elements []Renderable
}

func (g Group) Render() string {
	elements := make([]string, len(g.Elements))
	for i, elem := range g.Elements {
		elements[i] = elem.Render()
	}
	return fmt.Sprintf(`{"type":"group","elements":[%s]}`, strings.Join(elements, ","))
}

// AddToStack Helper functions for Io struct
func (io *Io) AddToStack(element Renderable) {
	io.stack = append(io.stack, element)
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

func (i *Input) Text(label string, options ...func(*TextInput)) Inputable {
	input := &TextInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	i.io.AddToStack(input)
	return input
}

func (i *Input) Boolean(label string, options ...func(*BooleanInput)) Inputable {
	input := &BooleanInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	i.io.AddToStack(input)
	return input
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
