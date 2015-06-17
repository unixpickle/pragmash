package pragmash

import "errors"

var (
	ErrEOFAfterLineContinuation = errors.New("EOF after line continuation")
	ErrMissingOpenCurlyBrace    = errors.New("missing open curly brace")
	ErrUnexpectedCloseParen     = errors.New("unexpected ')'")
	ErrEscapeCodeUnderflow      = errors.New("escape code is too short")
	ErrMissingEndQuote          = errors.New("missing string terminator")
)
