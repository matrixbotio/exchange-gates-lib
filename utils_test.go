package matrixgates

import (
	"strconv"
	"testing"
)

func TestGetFloatPrecision(t *testing.T) {
	floatVal := 56.13954
	precisionExpected := 5
	precision := GetFloatPrecision(floatVal)
	if precision != precisionExpected {
		t.Fatalf("failed to count float value precision. Received " +
			strconv.Itoa(precision) + ", expected " + strconv.Itoa(precisionExpected))
	}
}
