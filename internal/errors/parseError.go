package errors

import "fmt"

type ParseError struct {
	ParsedString string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parse error of string: %s", e.ParsedString)
}
