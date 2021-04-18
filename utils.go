package matrixgates

import "math"

//GetFloatPrecision returns the number of decimal places in a float
func GetFloatPrecision(f float64) int {
	return int(math.Ceil(math.Log10(math.Floor(1 / f))))
}
