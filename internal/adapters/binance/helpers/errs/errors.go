package errs

import "strings"

func IsErrorAboutUnknownOrder(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Unknown order sent")
}
