package parseargs

import "strings"

type State int

const (
	StateNormal = iota
	StateSingleQuote
	StateDoubleQuote
	StateEscape
	StateEscapeDoubleQuote
)

func TrimString(argin string) []string {
	var result []string
	var state State = StateNormal
	var buf strings.Builder
	quote := strings.ReplaceAll(argin, "''", "")
	arg := strings.ReplaceAll(quote, "\"\"", "")
	for _, c := range arg {
		switch state {
		case StateNormal:
			if c == '\n' || c == ' ' {
				if buf.Len() > 0 {
					result = append(result, buf.String())
					buf.Reset()
				}
			} else if c == '\\' {
				state = StateEscape
			} else if c == '"' {
				if buf.Len() > 0 {
					result = append(result, buf.String())
					buf.Reset()
				}
				state = StateDoubleQuote
			} else if c == '\'' {
				if buf.Len() > 0 {
					result = append(result, buf.String())
					buf.Reset()
				}
				state = StateSingleQuote
			} else {
				buf.WriteRune(c)
			}
		case StateSingleQuote:
			if c == '\'' {
				// if buf.Len() > 0 {
				// 	result = append(result, buf.String())
				// 	buf.Reset()
				// }
				state = StateNormal
			} else {
				buf.WriteRune(c)
			}
		case StateDoubleQuote:
			if c == '"' {
				// buf.WriteRune(c)
				state = StateNormal
				// if buf.Len() > 0 {
				// 	result = append(result, buf.String())
				// 	buf.Reset()
				// }
				// state = StateNormal
			} else if c == '\\' {
				state = StateEscapeDoubleQuote
			} else {
				buf.WriteRune(c)
			}
		case StateEscape:
			buf.WriteRune(c)
			state = StateNormal
		case StateEscapeDoubleQuote:
			buf.WriteRune(c)
			state = StateDoubleQuote
		}
	}
	if buf.Len() > 0 {
		result = append(result, buf.String())
	}
	return result
}
