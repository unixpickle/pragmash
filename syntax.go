package pragmash

import (
	"bytes"
	"errors"
)

var BlockInitiatingKeywords = []string{"if", "else", "while", "try", "for", "def"}

// A Token is either a nested command or a string literal.
type Token struct {
	// Command is nil if the token is a string literal.
	// Otherwise, it is an array of tokens in the nested command.
	Command []Token

	// Text is the token's text if it is a string literal.
	// If the token is a nested command, Text is "".
	Text string

	// Bare is true if the token was an unquoted string with no escapes.
	Bare bool
}

// A SyntaxLine is a logical line which has been parsed.
type SyntaxLine struct {
	// BlockOpen is true if the line ends with a { and begins with a block-initiating token.
	BlockOpen bool

	// BlockClose is true if the line begins with a }.
	BlockClose bool

	// Tokens stores the parsed tokens on the line, not including the curly braces that were taken
	// into account for BlockOpen and BlockClose.
	Tokens []Token

	// Number is a physical line number.
	Number int
}

// A SyntaxParser reads logical lines one at a time and parses them.
// It will ignore empty or commented lines.
type SyntaxParser struct {
	Reader LogicalLineReader
}

// ReadSyntaxLine reads and parses the next non-empty uncommented line.
// An error is returned if the underlying reader fails or if a syntax error is encountered.
func (s SyntaxParser) ReadSyntaxLine() (*SyntaxLine, error) {
	for {
		line, num, err := s.Reader.ReadLine()
		if err != nil {
			return nil, err
		} else if len(line) == 0 {
			continue
		} else if line[0] == '#' {
			continue
		}
		return parseLine(line, num)
	}
}

func parseLine(text string, num int) (*SyntaxLine, error) {
	line := &SyntaxLine{false, false, []Token{}, num}
	buffer := bytes.NewBufferString(text)
	for buffer.Len() > 0 {
		if token, err := readNextToken(buffer); err != nil {
			return nil, err
		} else {
			line.Tokens = append(line.Tokens, *token)
		}
	}
	return processCurlyBraces(line)
}

func processCurlyBraces(l *SyntaxLine) (*SyntaxLine, error) {
	if len(l.Tokens) == 0 {
		panic("there should always be tokens here")
	}
	if l.Tokens[0].Text == "}" && l.Tokens[0].Bare {
		l.Tokens = l.Tokens[1:]
		l.BlockClose = true
	}

	isOpenKeyword := false
	if len(l.Tokens) > 0 && l.Tokens[0].Bare {
		for _, keyword := range BlockInitiatingKeywords {
			if l.Tokens[0].Text == keyword {
				isOpenKeyword = true
				break
			}
		}
	}

	if isOpenKeyword {
		if l.Tokens[len(l.Tokens)-1].Text != "{" {
			return nil, ErrMissingOpenCurlyBrace
		}
		l.Tokens = l.Tokens[:len(l.Tokens)-1]
		l.BlockOpen = true
	}

	return l, nil
}

func readNextToken(buffer *bytes.Buffer) (*Token, error) {
	firstRune, _, err := buffer.ReadRune()
	if err != nil {
		return nil, err
	}
	switch firstRune {
	case '"':
		if str, err := readDoubleQuotedString(buffer); err != nil {
			return nil, err
		} else {
			return &Token{nil, str, false}, nil
		}
	case '\'':
		if str, err := readSingleQuotedString(buffer); err != nil {
			return nil, err
		} else {
			return &Token{nil, str, false}, nil
		}
	case '(':
		if tokens, err := readNestedCommand(buffer); err != nil {
			return nil, err
		} else {
			return &Token{tokens, "", false}, nil
		}
	case ')':
		return nil, ErrUnexpectedCloseParen
	default:
		buffer.UnreadRune()
		if token, err := readBareString(buffer); err != nil {
			return nil, err
		} else {
			return token, nil
		}
	}
}

func readDoubleQuotedString(buffer *bytes.Buffer) (string, error) {
	// TODO: this
	return "", errors.New("not yet implemented")
}

func readSingleQuotedString(buffer *bytes.Buffer) (string, error) {
	// TODO: this
	return "", errors.New("not yet implemented")
}

func readBareString(buffer *bytes.Buffer) (*Token, error) {
	// TODO: this
	return nil, errors.New("not yet implemented")
}

func readNestedCommand(buffer *bytes.Buffer) ([]Token, error) {
	// TODO: this
	return nil, errors.New("not yet implemented")
}

func readEscapeSequence(buffer *bytes.Buffer) (string, error) {
	// TODO: this
	return "", errors.New("not yet implemented")
}
