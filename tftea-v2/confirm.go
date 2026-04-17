// Copyright (C) 2025 T-Force I/O
//
// TFtea is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package tftea

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// ConfirmModel contains internal states of the Confirm.
type ConfirmModel struct {
	label     string
	value     bool
	cancelled bool
}

// Return a new Confirm instance with default config.
func NewConfirm() *ConfirmModel {
	return &ConfirmModel{
		label: "Please confirm:",
		value: false,
	}
}

// Set label.
func (m *ConfirmModel) WithLabel(label string) *ConfirmModel {
	m.label = label
	return m
}

// Set default value.
func (m *ConfirmModel) WithValue(value bool) *ConfirmModel {
	m.value = value
	return m
}

// Display the Confirm to the terminal
func (m *ConfirmModel) Run() (bool, error) {
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return false, err
	}
	final, ok := result.(ConfirmModel)
	if !ok {
		return false, ErrUnexpectedError
	}
	if final.cancelled {
		return false, ErrActionCancelled
	}
	return final.value, nil
}

// Bubbletea lifecycle implementation: Init
func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

// Bubbletea lifecycle implementation: Update
func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "y", "Y":
			m.value = true
			return m, tea.Quit
		case "n", "N":
			m.value = false
			return m, tea.Quit
		case "left", "right", "tab":
			m.value = !m.value
		case "enter":
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.value = false
			m.cancelled = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// Bubbletea lifecycle implementation: View
func (m ConfirmModel) View() tea.View {
	yes, no := "  Yes  ", "  No  "
	if m.value {
		yes = "[ Yes ]"
	} else {
		no = "[ No  ]"
	}
	s := fmt.Sprintf("%s\n\n  %s    %s\n", labelStyle.Render(m.label), yes, no)
	s += fmt.Sprintf("\n%s yes   %s no   %s navigate   %s confirm   %s cancel\n",
		keyStyle.Render(" Y "),
		keyStyle.Render(" N "),
		keyStyle.Render(" ←/→/tab "),
		keyStyle.Render(" enter "),
		keyStyle.Render(" esc "),
	)
	return tea.NewView(s)
}
