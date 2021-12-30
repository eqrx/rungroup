package rungroup

import "strings"

// Error wraps multiple errors into one error instance. It does not support unwrapping since the current interface
// design of Go allows only for one child of an error to be unwrapped. If you need to know the concrete types
// please go over Errs manually.
type Error struct {
	Errs []error
}

func (e Error) Error() string {
	errsLen := len(e.Errs)
	switch errsLen {
	case 0:
		panic("empty errs")
	case 1:
		return e.Errs[0].Error()
	default:
		strs := make([]string, 0, errsLen)
		for _, err := range e.Errs {
			strs = append(strs, err.Error())
		}

		return "multiple errors: [\"" + strings.Join(strs, "\", \"") + "\"]"
	}
}
