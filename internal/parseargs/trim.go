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
	Command         string
	Args            []string
	IsRedirect      bool
	IsRedirectError bool
	RedirectFile    string
}

func TrimString(argin string) Args {
	var result Args = Args{
		Command:         "",
		Args:            []string{},
		IsRedirect:      false,
		RedirectFile:    "",
		IsRedirectError: false,
	}
	if argin == "" {
		return result
	}
	var args []string
	var state State = StateNormal
	var buf strings.Builder

	arg := argin

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
				state = StateDoubleQuote
			} else if c == '\'' {
				state = StateSingleQuote
			} else {
				buf.WriteRune(c)
			}
		case StateSingleQuote:
			if c == '\'' {
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

	var countArgs int = len(args)
	for i, arg := range args[1:] {
		if result.IsRedirect || result.IsRedirectError {
			result.RedirectFile = arg
			break
		}
		if arg == ">" || arg == "1>" {
			result.IsRedirect = true
			countArgs = i + 1
		}
		if arg == "2>" {
			result.IsRedirectError = true
			countArgs = i + 1
		}
	}
	result.Command = args[0]
	result.Args = args[1:countArgs]
	return result
}
