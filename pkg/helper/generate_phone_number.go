package helper

import (
	"fmt"
	"math/rand"
)

func GeneratePhoneNumber() string {

	digits := make([]int, 10)
	for i := 0; i < 10; i++ {
		digits[i] = rand.Intn(10)
	}

	return fmt.Sprintf("+(%d%d%d) %d%d%d-%d%d%d%d",
		digits[0], digits[1], digits[2], digits[3], digits[4], digits[5], digits[6], digits[7], digits[8], digits[9])
}
