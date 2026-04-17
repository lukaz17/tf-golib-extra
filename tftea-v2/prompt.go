// Copyright (C) 2025 T-Force I/O
//
// TFtea is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package tftea

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Func to validate user input in the prompt.
type PromptValidateFunc func(value string) error

// PromptModel contains internal states of the Prompt.
type PromptModel struct {
	label         string
	value         string
	input         textinput.Model
	validateFuncs []PromptValidateFunc
	err           error
	cancelled     bool
}

// Return a new Prompt instance with default config.
func NewPrompt() *PromptModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.SetWidth(32)
	ti.Focus()

	return &PromptModel{
		label: "Please input:",
		input: ti,

		validateFuncs: make([]PromptValidateFunc, 0),
	}
}

// Set label.
func (m *PromptModel) WithLabel(label string) *PromptModel {
	m.label = label
	return m
}

// Set default value.
func (m *PromptModel) WithValue(value string) *PromptModel {
	m.value = value
	m.input.SetValue(value)
	return m
}

// Set placeholder text when the input is empty.
func (m *PromptModel) WithPlaceholder(placeholder string) *PromptModel {
	m.input.Placeholder = placeholder
	return m
}

// Add validation function enforce the rule of input value.
// Multiple validation functions are supported and follow the order they are added.
func (m *PromptModel) WithValidation(fn PromptValidateFunc) *PromptModel {
	m.validateFuncs = append(m.validateFuncs, fn)
	return m
}

// Set maximum number of characters allowed.
func (m *PromptModel) WithCharLimit(limit int) *PromptModel {
	m.input.CharLimit = limit
	return m
}

// Set visible width of the input field.
func (m *PromptModel) WithWidth(width int) *PromptModel {
	m.input.SetWidth(width)
	return m
}

// Display the Prompt to the terminal
func (m *PromptModel) Run() (string, error) {
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	final, ok := result.(PromptModel)
	if !ok {
		return "", ErrUnexpectedError
	}
	if final.cancelled {
		return "", ErrActionCancelled
	}
	return final.value, nil
}

// Bubbletea lifecycle implementation: Init
func (m PromptModel) Init() tea.Cmd {
	return textinput.Blink
}

// Bubbletea lifecycle implementation: Update
func (m PromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.err = nil
		switch msg.String() {
		case "enter":
			value := m.input.Value()
			for _, fn := range m.validateFuncs {
				if err := fn(value); err != nil {
					m.err = err
					return m, nil
				}
			}
			m.value = value
			m.err = nil
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// Bubbletea lifecycle implementation: View
func (m PromptModel) View() tea.View {
	s := fmt.Sprintf("%s\n\n%s\n", labelStyle.Render(m.label), m.input.View())
	if m.err != nil {
		s += fmt.Sprintf("\n  error: %s\n", m.err.Error())
	}
	s += fmt.Sprintf("\n%s confirm   %s cancel  \n",
		keyStyle.Render(" enter "),
		keyStyle.Render(" esc "),
	)
	return tea.NewView(s)
}
