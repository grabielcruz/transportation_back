package utility

import (
	"math"
	"math/rand"
	"time"
)

var seed = time.Now().Unix()

func GetRandomString(length int) string {
	rand.Seed(change_seed())

	ran_str := ""
	for i := 0; i < length; i++ {
		ran_str += string(rune(65 + rand.Intn(25)))
	}
	return ran_str
}

func GetRandomBalance() float64 {
	rand.Seed(change_seed())
	return math.Round(rand.Float64()*10000) / 100
}

func GetRandomBoolean() bool {
	rand.Seed(change_seed())
	return rand.Intn(2) == 1
}

func GetRandomCurrency() string {
	curr := "USD"
	if GetRandomBoolean() {
		curr = "VES"
	}
	return curr
}

func change_seed() int64 {
	seed = seed + 1
	return seed
}

func GenerateDummyData() DummyType {
	dummy := DummyType{}
	dummy.Hace = GetRandomString(4)
	dummy.Hola = GetRandomString(5)
	dummy.Que = GetRandomString(6)
	return dummy
}
