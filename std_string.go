package pragmash

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// StdString implements ways of manipulating or creating strings
type StdString struct{}

// Chars returns an array with each character from a string.
// Each newline character will be encoded as the two-character escape
// sequence "\\n".
func (_ StdString) Chars(s string) []string {
	runes := []rune(s)
	resArr := make([]string, len(runes))
	for i, x := range runes {
		if x == '\n' {
			resArr[i] = "\\n"
		} else {
			resArr[i] = string(x)
		}
	}
	return resArr
}

// Chr takes a list of numbers (representing bytes) and turns it into a string.
func (_ StdString) Chr(list []int) string {
	res := bytes.Buffer{}
	for _, x := range list {
		res.WriteByte(byte(x))
	}
	return res.String()
}

// Echo joins its arguments with spaces.
func (_ StdString) Echo(args ...string) string {
	return strings.Join(args, " ")
}

// Escape replaces backslashes with double-backslashes and newlines with "\n".
func (_ StdString) Escape(str string) string {
	s := strings.Replace(str, "\\", "\\\\", -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}

// HasPrefix returns true if the first argument begins with the second.
func (_ StdString) HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix returns true if the first argument ends with the second.
func (_ StdString) HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// IsDigit returns true if the provided argument is a number.
func (_ StdString) IsDigit(s string) bool {
	runes := []rune(s)
	if len(runes) != 1 {
		return false
	}
	return unicode.IsDigit(runes[0])
}

// IsLetter returns true if the provided argument is a capital or lowercase
// character in the English alphabet.
func (_ StdString) IsLetter(s string) bool {
	runes := []rune(s)
	if len(runes) != 1 {
		return false
	}
	return unicode.IsLetter(runes[0])
}

// Join joins its arguments without spaces.
func (_ StdString) Join(args ...string) string {
	var buffer bytes.Buffer
	for _, s := range args {
		buffer.WriteString(s)
	}
	return buffer.String()
}

// Len returns the length of a string in bytes.
func (_ StdInternal) Len(val string) int {
	return len(val)
}

// Lowercase joins its arguments with spaces and returns the result, converted
// to lower-case.
func (s StdString) Lowercase(args ...string) string {
	return strings.ToLower(s.Echo(args...))
}

// Match runs a regular expression on a string.
func (_ StdString) Match(expr, haystack string) ([]string, error) {
	// Evaluate the regular expression.
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	res := r.FindAllStringSubmatch(haystack, -1)

	// Return the result as a massive list.
	list := make([]string, 0)
	for _, x := range res {
		for _, y := range x {
			list = append(list, y)
		}
	}
	return list, nil
}

// Ord returns a list of numeric bytes given a string.
func (_ StdString) Ord(str string) []int {
	res := make([]int, len(str))
	for i := 0; i < len(str); i++ {
		res[i] = int(str[i])
	}
	return res
}

// PadZero pads a string with zeroes on the left until it's a certain length.
func (_ StdString) PadZero(length int, str string) string {
	for len(str) < length {
		str = "0" + str
	}
	return str
}

// Rep replaces all occurences of a string with another string.
func (_ StdString) Rep(s, old, replacement string) string {
	return strings.Replace(s, old, replacement, -1)
}

// Repreg replaces all occurences of a regular expression with an expandable
// expression.
func (_ StdString) Repreg(s, expr, replacement string) (string, error) {
	// Evaluate the regular expression.
	r, err := regexp.Compile(expr)
	if err != nil {
		return "", err
	}
	// Perform the replacement
	return r.ReplaceAllString(s, replacement), nil
}

// Substr returns a substring of a large string.
func (_ StdString) Substr(s string, start int, e ...int) (string, error) {
	if len(s) == 0 {
		return "", nil
	}

	var end int
	if len(e) > 1 {
		return "", errors.New("expected 2 or 3 arguments")
	} else if len(e) == 1 {
		end = e[0]
	} else {
		end = len(s)
	}

	// Any inputs are sanitized and accepted.
	if start < 0 {
		start = 0
	} else if start > len(s) {
		start = len(s)
	}
	if end < start {
		end = start
	} else if end > len(s) {
		end = len(s)
	}

	return s[start:end], nil
}

// Unescape replaces "\\" with "\" and "\n" with a newline.
func (_ StdString) Unescape(arg string) string {
	s := strings.Replace(arg, "\\n", "\n", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}

// Uppercase joins its arguments with spaces and returns the result, converted
// to upper-case.
func (s StdString) Uppercase(args ...string) string {
	return strings.ToUpper(s.Echo(args...))
}
