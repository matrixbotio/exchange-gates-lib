package matrixgates

import (
	"log"
	"math"
)

//GetFloatPrecision returns the number of decimal places in a float
func GetFloatPrecision(f float64) int {
	return int(math.Ceil(math.Log10(math.Floor(1 / f))))
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
