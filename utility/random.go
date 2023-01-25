package utility

import (
	"math/rand"
	"strings"
	"time"
)

func GetRandomNumbersInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandomString(n int) (str string) {
	rand.Seed(time.Now().UnixNano())

	var finalString = ""

	randgenS := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	s := strings.Split(randgenS, "")

	for j := 1; j <= n; j++ {
		randIdx := rand.Intn(len(s))
		finalString += s[randIdx]
	}

	return finalString
}
