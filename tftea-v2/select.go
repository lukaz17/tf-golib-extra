// Copyright (C) 2025 T-Force I/O
//
// TFtea is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package tftea

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
)

// SelectOption represents a key-value pair used as a selectable item in Select.
type SelectOption struct {
	Key   string
	Label string
}

// Return new SelectOption instance.
func NewSelectOption(key, label string) *SelectOption {
	return &SelectOption{Key: key, Label: label}
}

// Func to filter options in the Select.
type SelectFilterFunc func(query string, all []*SelectOption) []*SelectOption

// SelectModel contains internal states of the Select.
type SelectModel struct {
	label           string
	selected        map[string]bool
	filter          textinput.Model
	multiple        bool
	options         []*SelectOption
	filterFunc      SelectFilterFunc
	filtered        []*SelectOption
	minVisibleItems int
	maxVisibleItems int
	cursor          int
	pageSize        int
	pageOffset      int
	termHeight      int
	termWidth       int
	err             error
	cancelled       bool
}

// Return a new Select instance with default config.
func NewSelect() *SelectModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Placeholder = "type to filter..."
	ti.SetWidth(32)
	ti.Focus()

	return &SelectModel{
		label:           "Please select:",
		selected:        make(map[string]bool),
		filter:          ti,
		options:         make([]*SelectOption, 0),
		filterFunc:      selectFilter,
		minVisibleItems: 3,
		maxVisibleItems: 10,
	}
}

// Set label.
func (m *SelectModel) WithLabel(label string) *SelectModel {
	m.label = label
	return m
}

// Set default value.
func (m *SelectModel) WithValue(value string) *SelectModel {
	m.selected = make(map[string]bool)
	m.selected[value] = true
	return m
}

// Set default values.
func (m *SelectModel) WithValues(value []string) *SelectModel {
	m.selected = make(map[string]bool)
	for _, v := range value {
		m.selected[v] = true
	}
	return m
}

// Set multiple selection mode.
func (m *SelectModel) WithMultiSelect(multiple bool) *SelectModel {
	m.multiple = multiple
	return m
}

// Set initial options.
func (m *SelectModel) WithOptions(options []*SelectOption) *SelectModel {
	m.options = options
	return m
}

// Add one option to current options list.
func (m *SelectModel) WithExtraOption(key, label string) *SelectModel {
	opt := NewSelectOption(key, label)
	m.options = append(m.options, opt)
	return m
}

// Set filter function.
func (m *SelectModel) WithFilter(fn SelectFilterFunc) *SelectModel {
	m.filterFunc = fn
	return m
}

// Set maximum number of items per page.
func (m *SelectModel) WithMaxVisibleItems(n int) *SelectModel {
	m.maxVisibleItems = n
	return m
}

// Display the Prompt to the terminal
func (m *SelectModel) Run() ([]string, error) {
	m.updateTermSize()
	m.applyFilter()
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}
	final, ok := result.(SelectModel)
	if !ok {
		return nil, ErrUnexpectedError
	}
	if final.cancelled {
		return nil, ErrActionCancelled
	}

	return final.values(), nil
}

func (m *SelectModel) values() []string {
	v := make([]string, 0)
	for k, selected := range m.selected {
		if selected {
			v = append(v, k)
		}
	}
	return v
}

func (m *SelectModel) calcPageSize() int {
	overhead := 8
	available := m.termHeight - overhead
	if available < m.minVisibleItems {
		available = m.minVisibleItems
	}
	if available > m.maxVisibleItems {
		available = m.maxVisibleItems
	}
	return available
}

func (m *SelectModel) applyFilter() {
	m.filtered = m.filterFunc(m.filter.Value(), m.options)
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.clampPageOffset()
}

func (m *SelectModel) applyTermSize(w, h int) {
	m.termWidth = w
	m.termHeight = h
	fw := w - 6
	if fw < 10 {
		fw = 10
	}
	m.filter.SetWidth(fw)
	m.pageSize = m.calcPageSize()
}

func (m *SelectModel) updateTermSize() {
	w, h, err := term.GetSize(os.Stdout.Fd())
	if err != nil || w <= 0 || h <= 0 {
		w, _ = strconv.Atoi(os.Getenv("COLUMNS"))
		h, _ = strconv.Atoi(os.Getenv("LINES"))
	}
	if w <= 0 || h <= 0 {
		w, h = 80, 24
	}
	m.applyTermSize(w, h)
}

func (m *SelectModel) clampPageOffset() {
	if m.cursor < m.pageOffset {
		m.pageOffset = m.cursor
	}
	if m.cursor >= m.pageOffset+m.pageSize {
		m.pageOffset = m.cursor - m.pageSize + 1
	}
	if m.pageOffset < 0 {
		m.pageOffset = 0
	}
}

// Bubbletea lifecycle implementation: Init
func (m SelectModel) Init() tea.Cmd {
	return textinput.Blink
}

// Bubbletea lifecycle implementation: Update
func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.err = nil
		switch msg.String() {
		case "up":
			if len(m.filtered) == 0 {
				return m, nil
			}
			m.cursor = (m.cursor - 1 + len(m.filtered)) % len(m.filtered)
			if m.cursor == len(m.filtered)-1 {
				m.pageOffset = len(m.filtered) - m.pageSize
				if m.pageOffset < 0 {
					m.pageOffset = 0
				}
			}
			m.clampPageOffset()
			return m, nil
		case "down":
			if len(m.filtered) == 0 {
				return m, nil
			}
			m.cursor = (m.cursor + 1) % len(m.filtered)
			if m.cursor == 0 {
				m.pageOffset = 0
			}
			m.clampPageOffset()
			return m, nil
		case "left":
			if len(m.filtered) == 0 {
				return m, nil
			}
			m.cursor -= m.pageSize
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.pageOffset -= m.pageSize
			if m.pageOffset < 0 {
				m.pageOffset = 0
			}
			return m, nil
		case "right":
			last := len(m.filtered) - 1
			if last < 0 {
				return m, nil
			}
			m.cursor += m.pageSize
			if m.cursor > last {
				m.cursor = last
			}
			m.pageOffset += m.pageSize
			maxOffset := last - m.pageSize + 1
			if maxOffset < 0 {
				maxOffset = 0
			}
			if m.pageOffset > maxOffset {
				m.pageOffset = maxOffset
			}
			return m, nil
		case "tab":
			if m.multiple && len(m.filtered) > 0 {
				key := m.filtered[m.cursor].Key
				m.selected[key] = !m.selected[key]
			}
			return m, nil
		case "enter":
			if m.multiple {
				v := m.values()
				if len(v) == 0 {
					m.err = ErrNothingSelected
					return m, nil
				}
			} else {
				if len(m.filtered) == 0 {
					m.err = ErrNothingSelected
					return m, nil
				}
				cur := m.filtered[m.cursor]
				m.selected = map[string]bool{cur.Key: true}
			}
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.filter, cmd = m.filter.Update(msg)
		m.applyFilter()
		return m, cmd

	case tea.WindowSizeMsg:
		m.applyTermSize(msg.Width, msg.Height)
		m.clampPageOffset()
		return m, tea.ClearScreen
	}

	var cmd tea.Cmd
	m.filter, cmd = m.filter.Update(msg)
	return m, cmd
}

// Bubbletea lifecycle implementation: View
func (m SelectModel) View() tea.View {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n\n", labelStyle.Render(m.label)))
	sb.WriteString(fmt.Sprintf("Filter: %s\n\n", m.filter.View()))

	total := len(m.filtered)
	if total == 0 {
		sb.WriteString("  (No options)\n")
		if m.err != nil {
			sb.WriteString(fmt.Sprintf("  error: %s\n", m.err.Error()))
		}
	} else {
		end := m.pageOffset + m.pageSize
		if end > total {
			end = total
		}
		for i := m.pageOffset; i < end; i++ {
			opt := m.filtered[i]

			cursor := "  "
			if i == m.cursor {
				cursor = "> "
			}

			var tick string
			if m.multiple {
				if m.selected[opt.Key] {
					tick = "[x] "
				} else {
					tick = "[ ] "
				}
			}
			sb.WriteString(fmt.Sprintf("%s%s%s\n", cursor, tick, opt.Label))
		}

		if m.err != nil {
			sb.WriteString(fmt.Sprintf("\n  error: %s\n", m.err.Error()))
		} else {
			padding := len(strconv.Itoa(total))
			totalPages := (total + m.pageSize - 1) / m.pageSize
			paddingPage := len(strconv.Itoa(totalPages))
			currentPage := m.pageOffset/m.pageSize + 1
			sb.WriteString(fmt.Sprintf("\n  %*d–%*d of %*d  (page %*d/%*d)\n",
				padding, m.pageOffset+1,
				padding, end,
				padding, total,
				paddingPage, currentPage,
				paddingPage, totalPages,
			))
		}
	}

	if m.multiple {
		sb.WriteString(fmt.Sprintf("\n%s navigate   %s page   %s toggle   %s confirm   %s cancel  \n",
			keyStyle.Render(" ↑/↓ "),
			keyStyle.Render(" ←/→ "),
			keyStyle.Render(" tab "),
			keyStyle.Render(" enter "),
			keyStyle.Render(" esc "),
		))
	} else {
		sb.WriteString(fmt.Sprintf("\n%s navigate   %s page   %s confirm   %s cancel  \n",
			keyStyle.Render(" ↑/↓ "),
			keyStyle.Render(" ←/→ "),
			keyStyle.Render(" enter "),
			keyStyle.Render(" esc "),
		))
	}

	v := tea.NewView(sb.String())
	v.AltScreen = true
	return v
}

// Filter options case-insensitive on the option label.
func selectFilter(query string, all []*SelectOption) []*SelectOption {
	result := make([]*SelectOption, 0, len(all))
	if query == "" {
		return append(result, all...)
	}
	lower := strings.ToLower(query)
	for _, o := range all {
		if strings.Contains(strings.ToLower(o.Label), lower) {
			result = append(result, o)
		}
	}
	return result
}
