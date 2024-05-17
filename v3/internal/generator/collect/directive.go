package collect

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/wailsapp/wails/v3/internal/flags"
)

// IsDirective returns true if the given comment
// is a directive of the form //wails: + directive.
func IsDirective(comment string, directive string) bool {
	if strings.HasPrefix(comment, "//wails:"+directive) {
		length := len("//wails:") + len(directive)
		if len(comment) == length {
			return true
		}

		next, _ := utf8.DecodeRuneInString(comment[length:])
		return unicode.IsSpace(next)
	}

	return false
}

// ParseDirective extracts the argument portion of a //wails: + directive comment.
func ParseDirective(comment string, directive string) string {
	rawArg := comment[len("//wails:")+len(directive):]

	if directive != "inject" {
		return strings.TrimSpace(rawArg)
	}

	// wails:inject requires special parsing:
	// do not trim all surrounding space, just the one space
	// immediately after the directive name.
	_, wsize := utf8.DecodeRuneInString(rawArg)
	return rawArg[wsize:]
}

// ParseCondition parses an optional two-character condition prefix
// for include or inject directives.
// It returns the argument stripped of the prefix and the resulting condition.
// If the condition is malformed, ParseCondition returns a non-nil error.
func ParseCondition(argument string) (string, Condition, error) {
	cond, arg, present := strings.Cut(argument, ":")
	if !present {
		return cond, Condition{true, true, true, true}, nil
	}

	if len(cond) != 2 || !strings.ContainsRune("*jt", rune(cond[0])) || !strings.ContainsRune("*ci", rune(cond[1])) {
		return argument,
			Condition{true, true, true, true},
			fmt.Errorf("invalid condition code '%s': expected format is '(*|j|t)(*|c|i)'", cond)
	}

	condition := Condition{true, true, true, true}

	switch cond[0] {
	case 'j':
		condition.TS = false
	case 't':
		condition.JS = false
	}

	switch cond[1] {
	case 'c':
		condition.Interfaces = false
	case 'i':
		condition.Classes = false
	}

	return arg, condition, nil
}

type Condition struct {
	JS         bool
	TS         bool
	Classes    bool
	Interfaces bool
}

// Satisfied returns true when the condition described by the receiver
// is satisfied by the given configuration.
func (cond Condition) Satisfied(options *flags.GenerateBindingsOptions) bool {
	if options.TS {
		if options.UseInterfaces {
			return cond.TS && cond.Interfaces
		} else {
			return cond.TS && cond.Classes
		}
	} else {
		if options.UseInterfaces {
			return cond.JS && cond.Interfaces
		} else {
			return cond.JS && cond.Classes
		}
	}
}
