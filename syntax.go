package pragmash

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"unicode"
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
		if res, err := readSpace(buffer); err != nil {
			return nil, err
		} else if !res {
			return nil, ErrMissingWhitespace
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
		if str, err := readQuotedString(buffer, '"'); err != nil {
			return nil, err
		} else {
			return &Token{nil, str, false}, nil
		}
	case '\'':
		if str, err := readQuotedString(buffer, '\''); err != nil {
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

func readNestedCommand(buffer *bytes.Buffer) ([]Token, error) {
	tokens := []Token{}
	readSpace(buffer)
	for {
		if token, err := readNextToken(buffer); err == ErrUnexpectedCloseParen {
			break
		} else if err == io.EOF {
			return nil, ErrMissingCloseParen
		} else if err != nil {
			return nil, err
		} else {
			tokens = append(tokens, *token)
		}

		// We must check for a ')' before we attempt to read whitespace, since there needn't be any
		// whitespace before the ')'.
		if r, _, err := buffer.ReadRune(); err == nil {
			if r == ')' {
				break
			}
			buffer.UnreadRune()
		}

		if res, err := readSpace(buffer); err != nil {
			return nil, err
		} else if !res {
			return nil, ErrMissingWhitespace
		}
	}
	if len(tokens) == 0 {
		return nil, ErrEmptyParens
	}
	return tokens, nil
}

func readQuotedString(buffer *bytes.Buffer, quote rune) (string, error) {
	str := &bytes.Buffer{}
	for {
		rune, _, err := buffer.ReadRune()
		if err == io.EOF {
			return "", ErrMissingEndQuote
		} else if err != nil {
			return "", err
		} else if rune == quote {
			break
		} else if rune == '\\' {
			if seq, err := readEscapeSequence(buffer); err != nil {
				return "", err
			} else {
				if _, err := str.WriteRune(seq); err != nil {
					return "", err
				}
			}
		} else {
			if _, err := str.WriteRune(rune); err != nil {
				return "", err
			}
		}
	}
	return str.String(), nil
}

func readBareString(buffer *bytes.Buffer) (*Token, error) {
	str := &bytes.Buffer{}
	bare := true
	for buffer.Len() > 0 {
		rune, _, err := buffer.ReadRune()
		if err != nil {
			return nil, err
		} else if unicode.IsSpace(rune) || rune == ')' {
			buffer.UnreadRune()
			break
		} else if rune == '\\' {
			if seq, err := readEscapeSequence(buffer); err != nil {
				return nil, err
			} else {
				if _, err := str.WriteRune(seq); err != nil {
					return nil, err
				}
				bare = false
			}
		} else {
			if _, err := str.WriteRune(rune); err != nil {
				return nil, err
			}
		}
	}
	return &Token{nil, str.String(), bare}, nil
}

func readEscapeSequence(buffer *bytes.Buffer) (rune, error) {
	firstRune, _, err := buffer.ReadRune()
	if err != nil {
		return 0, ErrEscapeCodeUnderflow
	}
	switch firstRune {
	case '(', ')', '?', '\'', '"', ' ', '\\':
		return firstRune, nil
	case 'a':
		return '\a', nil
	case 'b':
		return '\b', nil
	case 'f':
		return '\f', nil
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil
	case 'v':
		return '\v', nil
	case 'x':
		return readNumericEscape(buffer, 2)
	case 'u':
		return readNumericEscape(buffer, 4)
	case 'U':
		return readNumericEscape(buffer, 8)
	default:
		if !unicode.IsDigit(firstRune) || firstRune == '8' || firstRune == '9' {
			break
		}
		buffer.UnreadRune()
		return readOctalEscape(buffer)
	}
	return 0, errors.New("invalid escape character: " + string(firstRune))
}

func readNumericEscape(b *bytes.Buffer, charCount int) (rune, error) {
	runes := make([]rune, 0, charCount)
	for i := 0; i < charCount; i++ {
		if r, _, err := b.ReadRune(); err != nil {
			if err == io.EOF {
				return 0, ErrEscapeCodeUnderflow
			}
			return 0, err
		} else {
			runes = append(runes, r)
		}
	}
	str := string(runes)
	if res, err := strconv.ParseUint(str, 16, charCount*4); err != nil {
		return 0, err
	} else {
		return rune(res), nil
	}
}

func readOctalEscape(b *bytes.Buffer) (rune, error) {
	runes := make([]rune, 0, 3)
	for i := 0; i < 3; i++ {
		if r, _, err := b.ReadRune(); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		} else if r >= '0' && r < '8' {
			runes = append(runes, r)
		} else {
			b.UnreadRune()
			break
		}
	}
	if len(runes) == 0 {
		return 0, ErrEscapeCodeUnderflow
	}
	str := string(runes)
	if res, err := strconv.ParseUint(str, 8, 8); err != nil {
		return 0, err
	} else {
		return rune(res), nil
	}
}

func readSpace(b *bytes.Buffer) (bool, error) {
	if b.Len() == 0 {
		return true, nil
	}
	gotSpace := false
	for b.Len() > 0 {
		if r, _, err := b.ReadRune(); err != nil {
			return false, err
		} else if unicode.IsSpace(r) {
			gotSpace = true
		} else {
			b.UnreadRune()
			break
		}
	}
	return gotSpace, nil
}
