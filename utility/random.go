package utility

import (
	"math/rand"
	"time"

	"github.com/gofrs/uuid"
)

func GetRandomNumbersInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandomString(length int) string {
	u, _ := uuid.NewV4()
	uuidStr := u.String()
	return (uuidStr + uuidStr[:length%36])[:length]
}
