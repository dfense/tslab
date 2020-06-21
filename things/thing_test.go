package things

import (
	"errors"
	"testing"
)

var (
	errFloatOutOfRange = errors.New("random generated float out of range")
)

// TestMinMaxFloat sample run 10000 times to make sure min/max borders were set correctly, and to show general usage
func TestMinMaxFloat(t *testing.T) {

	// run 1,000,000 times
	max := 100.00
	min := -50.00
	cnt := 0
	for cnt < 1000000 {

		floatValue := RFloat(min, max)
		if floatValue < min || floatValue > max {
			t.Errorf("%s: %f", errFloatOutOfRange, floatValue)
		}
		cnt++
	}
}
