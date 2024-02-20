package helper

import (
	"math/rand"
)

func GenerateRandomPrice(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
