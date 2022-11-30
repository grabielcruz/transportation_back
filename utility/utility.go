package utility

import (
	"math"
	"math/rand"
	"time"
)

func GetRandomString(length int) string {
	rand.Seed(time.Now().Unix())
	ran_str := ""
	for i := 0; i < length; i++ {
		ran_str += string(rune(65 + rand.Intn(25)))
	}
	return ran_str
}

func GetRandomBalance() float64 {
	return math.Round(rand.Float64()*10000) / 100
}

func GetRandomBoolean() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func GetRandomCurrency() string {
	curr := "USD"
	if GetRandomBoolean() {
		curr = "VES"
	}
	return curr
}
