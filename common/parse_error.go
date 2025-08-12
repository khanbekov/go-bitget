package common

import "fmt"

// Ошибки парсинга
type ParseError struct {
	Value  interface{} // Исходное значение
	Target string      // Целевой тип
	Msg    string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parsing failed: %v (type: %T) -> %s: %s", e.Value, e.Value, e.Target, e.Msg)
}
