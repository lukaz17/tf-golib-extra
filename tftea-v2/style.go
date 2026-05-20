// Copyright (C) 2025 T-Force I/O
//
// TFtea is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package tftea

import "charm.land/lipgloss/v2"

var (
	focusedOptionStyle = lipgloss.NewStyle().Reverse(true)
	labelStyle         = lipgloss.NewStyle().Bold(true)
	optionStyle        = lipgloss.NewStyle()
	shortcutStyle      = lipgloss.NewStyle().Bold(true).Reverse(true)
)
