// Copyright (C) 2025 T-Force I/O
//
// TFtea is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package tftea

import "errors"

var (
	ErrActionCancelled = errors.New("action cancelled")
	ErrUnexpectedError = errors.New("unexpected error")
)
