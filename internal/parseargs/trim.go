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

type Args struct {
	Command      string
	Args         []string
	IsRedirect   bool
	RedirectFile string
}

func TrimString(argin string) Args {
	var result Args = Args{
		Command:      "",
		Args:         []string{},
		IsRedirect:   false,
		RedirectFile: "",
	}
	if argin == "" {
		return result
	}
	var args []string
	var state State = StateNormal
	var buf strings.Builder

	arg := argin

	// quote := strings.ReplaceAll(argin, "''", "")
	//arg := strings.ReplaceAll(quote, "\"\"", "")
	for _, c := range arg {
		switch state {
		case StateNormal:
			if c == '\n' || c == ' ' {
				if buf.Len() > 0 {
					args = append(args, buf.String())
					buf.Reset()
				}
			} else if c == '\\' {
				state = StateEscape
			} else if c == '"' {
				// if buf.Len() > 0 {
				// 	args = append(args, buf.String())
				// 	buf.Reset()
				// }
				state = StateDoubleQuote
			} else if c == '\'' {
				// if buf.Len() > 0 {
				// 	args = append(args, buf.String())
				// 	buf.Reset()
				// }
				state = StateSingleQuote
			} else {
				buf.WriteRune(c)
			}
		case StateSingleQuote:
			if c == '\'' {
				// state = StateNormal
				// if buf.Len() > 0 {
				// 	args = append(args, buf.String())
				// 	buf.Reset()
				// }
				state = StateNormal
			} else {
				buf.WriteRune(c)
			}
		case StateDoubleQuote:
			if c == '"' {
				state = StateNormal
			} else if c == '\\' {
				state = StateEscapeDoubleQuote
			} else {
				buf.WriteRune(c)
			}
		case StateEscape:
			buf.WriteRune(c)
			state = StateNormal
		case StateEscapeDoubleQuote:
			if c == '"' {
				buf.WriteRune('"')
			} else if c == '\\' {
				buf.WriteRune('\\')
			} else {
				buf.WriteRune('\\')
				buf.WriteRune(c)
			}
			state = StateDoubleQuote
		}
	}
	if buf.Len() > 0 {
		args = append(args, buf.String())
	}

	var isRedirect bool = false
	var redirectFile string = ""
	var countArgs int = len(args)
	for i, arg := range args[1:] {
		if isRedirect {
			redirectFile = arg
			break
		}
		if arg == ">" || arg == "1>" {
			isRedirect = true
			countArgs = i
		}
	}
	result.Command = args[0]
	result.Args = args[1:countArgs]
	result.IsRedirect = isRedirect
	result.RedirectFile = redirectFile
	return result
}
