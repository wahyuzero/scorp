package collectors

import (
	"math"
)

func Round1(v float64) float64 {
	return math.Round(v*10) / 10
}
