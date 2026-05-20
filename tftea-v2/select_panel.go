// Copyright (C) 2025 T-Force I/O
//
// TFtea is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package tftea

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/tforce-io/tf-golib/opx"
)

// SelectPanel contains internal state of the SelectPanel.
type SelectPanelModel struct {
	label           string
	selectedSelect  int
	selectedOptions map[int]int
	selectLabels    map[int]string
	selectWidth     int
	optionLabels    map[int][]string
	optionWidth     int
	total           int
	hotkeyEnabled   bool
	hotkeys         []string
	err             error
	cancelled       bool
	confirmed       bool
}

// Return a new SelectPanel instance with default config.
func NewSelectPanel() *SelectPanelModel {
	return &SelectPanelModel{
		selectedOptions: make(map[int]int),
		selectLabels:    make(map[int]string),
		optionLabels:    make(map[int][]string),
		hotkeys:         []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"},
	}
}

// Set label.
func (m *SelectPanelModel) WithLabel(label string) *SelectPanelModel {
	m.label = label
	return m
}

// Set minimum width for select option.
func (m *SelectPanelModel) WithSelected(values []int) *SelectPanelModel {
	for i, v := range values {
		m.selectedOptions[i] = opx.Ternary(v >= 0, v, 0)
	}
	return m
}

// Add Select and its Options.
func (m *SelectPanelModel) WithSelect(label string, options []string) *SelectPanelModel {
	if len(options) == 0 {
		return m
	}
	m.selectLabels[m.total] = label
	m.optionLabels[m.total] = options
	m.total++
	return m
}

// Set enable for fast option switching.
func (m *SelectPanelModel) WithHotkey(enabled bool) *SelectPanelModel {
	m.hotkeyEnabled = enabled
	return m
}

// Set minimum width for select label.
func (m *SelectPanelModel) WithSelectWidth(width int) *SelectPanelModel {
	m.selectWidth = width
	return m
}

// Set minimum width for select option.
func (m *SelectPanelModel) WithOptionWidth(width int) *SelectPanelModel {
	m.optionWidth = width
	return m
}

// Bubbletea lifecycle implementation: Init
func (m SelectPanelModel) Init() tea.Cmd {
	return nil
}

// Display the SelectPanel to the terminal.
func (m *SelectPanelModel) Run() (map[int]int, error) {
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}
	final, ok := result.(SelectPanelModel)
	if !ok {
		return nil, ErrUnexpectedError
	}
	if final.cancelled {
		return nil, ErrActionCancelled
	}
	return final.selectedOptions, nil
}

// Bubbletea lifecycle implementation: Update.
func (m SelectPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m = m.switchSelect(false)
		case "down":
			m = m.switchSelect(true)
		case "left":
			m = m.switchCurrentOption(false)
		case "right":
			m = m.switchCurrentOption(true)
		case "1":
			m = m.switchOption(0, true)
		case "2":
			m = m.switchOption(1, true)
		case "3":
			m = m.switchOption(2, true)
		case "4":
			m = m.switchOption(3, true)
		case "5":
			m = m.switchOption(4, true)
		case "6":
			m = m.switchOption(5, true)
		case "7":
			m = m.switchOption(6, true)
		case "8":
			m = m.switchOption(7, true)
		case "9":
			m = m.switchOption(8, true)
		case "a", "A":
			m = m.switchOption(9, true)
		case "b", "B":
			m = m.switchOption(10, true)
		case "c", "C":
			m = m.switchOption(11, true)
		case "d", "D":
			m = m.switchOption(12, true)
		case "e", "E":
			m = m.switchOption(13, true)
		case "f", "F":
			m = m.switchOption(14, true)
		case "g", "G":
			m = m.switchOption(15, true)
		case "h", "H":
			m = m.switchOption(16, true)
		case "i", "I":
			m = m.switchOption(17, true)
		case "j", "J":
			m = m.switchOption(18, true)
		case "k", "K":
			m = m.switchOption(19, true)
		case "l", "L":
			m = m.switchOption(20, true)
		case "m", "M":
			m = m.switchOption(21, true)
		case "n", "N":
			m = m.switchOption(22, true)
		case "o", "O":
			m = m.switchOption(23, true)
		case "p", "P":
			m = m.switchOption(24, true)
		case "q", "Q":
			m = m.switchOption(25, true)
		case "r", "R":
			m = m.switchOption(26, true)
		case "s", "S":
			m = m.switchOption(27, true)
		case "t", "T":
			m = m.switchOption(28, true)
		case "u", "U":
			m = m.switchOption(29, true)
		case "v", "V":
			m = m.switchOption(30, true)
		case "w", "W":
			m = m.switchOption(31, true)
		case "x", "X":
			m = m.switchOption(32, true)
		case "y", "Y":
			m = m.switchOption(33, true)
		case "z", "Z":
			m = m.switchOption(34, true)
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			m.confirmed = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// Bubbletea lifecycle implementation: View
func (m SelectPanelModel) View() tea.View {
	var sb strings.Builder

	sb.WriteString(labelStyle.Render(m.label))
	sb.WriteString("\n\n")

	for i := 0; i < m.total; i++ {
		selectLabel := m.selectLabels[i]
		prefix := "  "
		lStyle := optionStyle

		if i == m.selectedSelect {
			prefix = "> "
			lStyle = focusedOptionStyle
		}

		hotkey := ""
		if m.hotkeyEnabled {
			hotkey = opx.Ternary(i < len(m.hotkeys), m.hotkeys[i]+". ", "   ")
		}
		paddedLabel := fmt.Sprintf("%-*s", m.selectWidth, fmt.Sprintf("%s%s ", hotkey, selectLabel))
		optionLabel := m.getLabel(i, m.selectedOptions[i])
		paddedValue := m.padString(optionLabel, m.optionWidth, true)

		sb.WriteString(fmt.Sprintf("%s%s < %s >\n",
			prefix,
			paddedLabel,
			lStyle.Render(paddedValue),
		))
	}

	sb.WriteString(fmt.Sprintf("\n%s navigate  %s change  %s change by key  %s confirm  %s cancel\n",
		shortcutStyle.Render(" ↑/↓ "),
		shortcutStyle.Render(" ←/→ "),
		shortcutStyle.Render(" 1-9/A-Z "),
		shortcutStyle.Render(" enter "),
		shortcutStyle.Render(" esc "),
	))

	return tea.NewView(sb.String())
}

func (m *SelectPanelModel) getLabel(selectIndex, optionIndex int) string {
	if selectIndex >= m.total || selectIndex < 0 {
		return ""
	}
	if optionIndex >= len(m.optionLabels[selectIndex]) || optionIndex < 0 {
		return ""
	}
	return m.optionLabels[selectIndex][optionIndex]
}

func (m *SelectPanelModel) padString(s string, width int, center bool) string {
	if len(s) >= width {
		return s
	}
	total := width - len(s)
	left := 1
	if center {
		left = total / 2
	}
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func (m SelectPanelModel) switchSelect(forward bool) SelectPanelModel {
	if forward {
		m.selectedSelect = (m.selectedSelect + 1) % m.total
	} else {
		m.selectedSelect = (m.selectedSelect - 1 + m.total) % m.total
	}
	return m
}

func (m SelectPanelModel) switchCurrentOption(forward bool) SelectPanelModel {
	optionCount := len(m.optionLabels[m.selectedSelect])
	offset := opx.Ternary(forward, 1, -1)
	m.selectedOptions[m.selectedSelect] = (m.selectedOptions[m.selectedSelect] + offset + optionCount) % optionCount
	return m
}

func (m SelectPanelModel) switchOption(selectIndex int, forward bool) SelectPanelModel {
	if (selectIndex >= m.total || selectIndex < 0) && m.hotkeyEnabled {
		return m
	}
	m.selectedSelect = selectIndex
	return m.switchCurrentOption(forward)
}
