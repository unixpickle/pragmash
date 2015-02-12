package pragmash

import (
	"errors"
	"strconv"
)

// OperatorRewrites can be used for a ReflectRunner's command rewrite table to
// replace the named operators with symbolic ones like *, +, [], <, etc.
var OperatorRewrites = map[string]string{
	"+": "add", "/": "div", "*": "mul", "-": "sub", "**": "pow",
	"[]": "subscript", "<=": "le", ">=": "ge", "<": "lt", ">": "gt", "=": "eq",
}

// StdOps implements the standard operators that are not implemented in StdMath.
type StdOps struct{}

// Eq implements the equality operator.
func (s StdOps) Eq(s1, s2 string) Value {
	return BoolValue(s1 == s2)
}

// Ge implements the >= operator.
func (s StdOps) Ge(n1, n2 Number) Value {
	return BoolValue(CompareNumbers(n1, n2) >= 0)
}

// Gt implements the > operator.
func (s StdOps) Gt(n1, n2 Number) Value {
	return BoolValue(CompareNumbers(n1, n2) > 0)
}

// Le implements the <= operator.
func (s StdOps) Le(n1, n2 Number) Value {
	return BoolValue(CompareNumbers(n1, n2) <= 0)
}

// Lt implements the < operator.
func (s StdOps) Lt(n1, n2 Number) Value {
	return BoolValue(CompareNumbers(n1, n2) < 0)
}

// Subscript gets a term from a list.
func (s StdOps) Subscript(strings []string, index int) (Value, error) {
	if index < 0 || index >= len(strings) {
		return nil, errors.New("subscript out of bounds: " +
			strconv.Itoa(index))
	}
	return StringValue(strings[index]), nil
}
