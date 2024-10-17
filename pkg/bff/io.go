package bff

import (
	"fmt"
	"strconv"
	"time"
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
// - input.number requests a number value
// - input.email requests an email value
// - input.slider requests a number value within a range
// - input.date requests a date value
// - input.textArea requests a text area value

// TODO:
// - input.richText requests a rich text value
// - input.url requests a URL value
// - input.time requests a date with time value
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
// - display.image Displays an image to the action user. One of url or buffer must be provided.
// - display.link Displays a button-styled action link to the action user. Can link to an external URL or to another action.
// - display.metadata Displays a series of label/value pairs in a variety of layout options.
// - display.code Displays a block of code to the action user.
// - display.html Displays rendered HTML to the action user.

// TODO:
// - display.grid  Displays data in a grid layout https://interval.com/docs/io-methods/display-grid
// - display.object Displays an object of nested data to the action user.
// - display.table Displays tabular data.
// - display.video Displays a video to the action user. One of url or buffer must be provided.

type Image struct {
	Url  string `json:"url,omitempty"`
	Alt  string `json:"alt,omitempty"`
	Size string `json:"size,omitempty"`
}

func (c Image) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "image", Data: c}
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

// LinkDisplay represents a button-styled action link
type LinkDisplay struct {
	Text string `json:"text"`
	Url  string `json:"url"`
	Type string `json:"type,omitempty"` // "default", "primary", "danger", etc.
}

func (l LinkDisplay) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "link", Data: l}
	return nil, nil
}

// HtmlDisplay represents rendered HTML content
type HtmlDisplay struct {
	Content string `json:"content"`
}

func (h HtmlDisplay) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "html", Data: h}
	return nil, nil
}

// CodeDisplay represents a block of code (already implemented, shown here for completeness)
type CodeDisplay struct {
	Code     string `json:"code"`
	Language string `json:"language,omitempty"`
}

func (c CodeDisplay) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "code", Data: c}
	return nil, nil
}

// MetadataItem represents a single label/value pair in the metadata display
type MetadataItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// MetadataDisplay represents a series of label/value pairs
type MetadataDisplay struct {
	Items  []MetadataItem `json:"items"`
	Layout string         `json:"layout,omitempty"` // "default", "card", "table", etc.
}

func (m MetadataDisplay) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "metadata", Data: m}
	return nil, nil
}

func (d *Display) Link(text string, url string, options ...func(*LinkDisplay)) {
	link := &LinkDisplay{Text: text, Url: url}
	for _, option := range options {
		option(link)
	}
	d.io.AddToStack(link)
}

func (d *Display) Html(content string) {
	d.io.AddToStack(HtmlDisplay{Content: content})
}

func (d *Display) Metadata(items []MetadataItem, options ...func(*MetadataDisplay)) {
	metadata := &MetadataDisplay{Items: items}
	for _, option := range options {
		option(metadata)
	}
	d.io.AddToStack(metadata)
}

// Option functions for customization

func WithLinkType(linkType string) func(*LinkDisplay) {
	return func(l *LinkDisplay) {
		l.Type = linkType
	}
}

func WithMetadataLayout(layout string) func(*MetadataDisplay) {
	return func(m *MetadataDisplay) {
		m.Layout = layout
	}
}

// Group combines multiple I/O method calls
type Group struct {
	Elements []Executable
}

func (g Group) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "group", Data: g}
	return nil, nil
}

// AddToStack adds the element to the stack and executes it -- returning the result of the execution
func (io *Io) AddToStack(element Executable) (any, error) {
	io.stack = append(io.stack, element)
	return element.Execute(io.input, io.output)
}

func (d *Display) Group(elements ...Executable) {
	d.io.AddToStack(Group{Elements: elements})
}
func (d *Display) Image(url string, alt string, size string) {
	d.io.AddToStack(Image{Url: url, Alt: alt, Size: size})
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

func (i *Input) Number(label string, options ...func(*NumberInput)) (int, error) {
	input := &NumberInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(v.(string))
	if err != nil {
		return 0, fmt.Errorf("expected number, got %T: %w", v, err)
	}
	return n, nil
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

// EmailInput represents an email input field
type EmailInput struct {
	InputBase
}

// SliderInput represents a slider input for number values within a range
type SliderInput struct {
	InputBase
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Step float64 `json:"step,omitempty"`
}

// DateInput represents a date input field
type DateInput struct {
	InputBase
	Min string `json:"min,omitempty"` // ISO 8601 date format (YYYY-MM-DD)
	Max string `json:"max,omitempty"` // ISO 8601 date format (YYYY-MM-DD)
}

// RichTextInput represents a rich text input field
type RichTextInput struct {
	InputBase
	InitialValue string `json:"initialValue,omitempty"`
}
type TextAreaInput struct {
	InputBase
	InitialValue string `json:"initialValue,omitempty"`
}

// URLInput represents a URL input field
type URLInput struct {
	InputBase
}

// TimeInput represents a time input field
type TimeInput struct {
	InputBase
	Min string `json:"min,omitempty"` // HH:mm format
	Max string `json:"max,omitempty"` // HH:mm format
}

// FileInput represents a file input field
type FileInput struct {
	InputBase
	Accept   string `json:"accept,omitempty"` // MIME types or file extensions
	Multiple bool   `json:"multiple,omitempty"`
}

// Implement Execute method for each new input type
func (e *EmailInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "emailInput", Data: e}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

func (s *SliderInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "sliderInput", Data: s}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

func (d *DateInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "dateInput", Data: d}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

func (r *RichTextInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "richTextInput", Data: r}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

func (r *TextAreaInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "textAreaInput", Data: r}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

func (u *URLInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "urlInput", Data: u}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

func (t *TimeInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "timeInput", Data: t}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

func (f *FileInput) Execute(input <-chan Message, output chan<- Message) (any, error) {
	output <- Message{Type: "fileInput", Data: f}
	m := <-input
	if m.Type != "input" {
		return nil, fmt.Errorf("expected input, got %s", m.Type)
	}
	return m.Data, nil
}

// Add new methods to the Input struct
func (i *Input) Email(label string, options ...func(*EmailInput)) (string, error) {
	input := &EmailInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (i *Input) Slider(label string, min, max float64, options ...func(*SliderInput)) (float64, error) {
	input := &SliderInput{InputBase: InputBase{Label: label}, Min: min, Max: max}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return 0, err
	}
	return v.(float64), nil
}

func (i *Input) Date(label string, options ...func(*DateInput)) (time.Time, error) {
	input := &DateInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("2006-01-02", v.(string))
}

func (i *Input) RichText(label string, options ...func(*RichTextInput)) (string, error) {
	input := &RichTextInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (i *Input) TextArea(label string, options ...func(*TextAreaInput)) (string, error) {
	input := &TextAreaInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (i *Input) URL(label string, options ...func(*URLInput)) (string, error) {
	input := &URLInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (i *Input) Time(label string, options ...func(*TimeInput)) (time.Time, error) {
	input := &TimeInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return time.Time{}, err
	}
	// parse
	return time.Parse("15:04", v.(string))
}

func (i *Input) File(label string, options ...func(*FileInput)) ([]string, error) {
	input := &FileInput{InputBase: InputBase{Label: label}}
	for _, option := range options {
		option(input)
	}
	v, err := i.io.AddToStack(input)
	if err != nil {
		return nil, err
	}
	arr, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("expected []string, got %T", v)
	}
	files := make([]string, 0, len(arr))
	for _, r := range arr {
		str, ok := r.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", r)
		}
		files = append(files, str)
	}

	return files, nil
}

// Option functions for customization
func WithStep(step float64) func(*SliderInput) {
	return func(s *SliderInput) {
		s.Step = step
	}
}

func WithDateRange(min, max string) func(*DateInput) {
	return func(d *DateInput) {
		d.Min = min
		d.Max = max
	}
}

func WithInitialValue(value string) func(*RichTextInput) {
	return func(r *RichTextInput) {
		r.InitialValue = value
	}
}

func WithTimeRange(min, max string) func(*TimeInput) {
	return func(t *TimeInput) {
		t.Min = min
		t.Max = max
	}
}

func WithAccept(accept string) func(*FileInput) {
	return func(f *FileInput) {
		f.Accept = accept
	}
}

func WithMultiple(multiple bool) func(*FileInput) {
	return func(f *FileInput) {
		f.Multiple = multiple
	}
}
