package utility

import (
	"math"
	"math/rand"
)

var i = new(int64)

func GetRandomString(length int) string {
	rand.Seed(increment_i())

	ran_str := ""
	for i := 0; i < length; i++ {
		ran_str += string(rune(65 + rand.Intn(25)))
	}
	return ran_str
}

func GetRandomBalance() float64 {
	rand.Seed(increment_i())
	return math.Round(rand.Float64()*10000) / 100
}

func GetRandomBoolean() bool {
	rand.Seed(increment_i())
	return rand.Intn(2) == 1
}

func GetRandomCurrency() string {
	curr := "USD"
	if GetRandomBoolean() {
		curr = "VES"
	}
	return curr
}

func increment_i() int64 {
	*i++
	return *i
}
