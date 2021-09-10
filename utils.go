package matrixgates

import (
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-stack/stack"
)

// GetFloatPrecision returns the number of decimal places in a float
func GetFloatPrecision(value float64) int {
	// if you put 15, then the test will fall,
	// because Float is rounded incorrectly
	maxPrecision := 14
	valueFormated := strconv.FormatFloat(math.Abs(value), 'f', maxPrecision, 64)
	valueParts := strings.Split(valueFormated, ".")
	if len(valueParts) <= 1 {
		return 0
	}
	log.Println(valueParts[1])
	valueLastPartTrimmed := strings.TrimRight(valueParts[1], "0")
	return len(valueLastPartTrimmed)
}

// LogNotNilError logs an array of errors and returns true if an error is found
func LogNotNilError(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			return true
		}
	}
	return false
}

// GetTrace - get stack string
func GetTrace() string {
	stackTrace := stack.Trace()
	if stackTrace == nil || len(stackTrace) == 0 {
		return ""
	}
	return stack.Trace().TrimRuntime().String()
}
