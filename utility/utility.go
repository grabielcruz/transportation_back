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
	balance := rand.Float64()
	return RoundToTwoDecimalPlaces(balance)
}

func getRandomFloat64() float64 {
	rand.Seed(change_seed())
	return rand.Float64()
}

func RoundToTwoDecimalPlaces(f float64) float64 {
	return math.Round(f*100) / 100
}

func GetRandomBoolean() bool {
	rand.Seed(change_seed())
	return rand.Intn(2) == 1
}

func GetRandomCurrency() string {
	curr := "USD"
	if GetRandomBoolean() {
		curr = "VED"
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

func GetSliceOfAmounts(total int) []float64 {
	nums := []float64{}
	for i := 1; i < total+1; i++ {
		f := getRandomFloat64() * 10
		if i%10 == 0 {
			f *= -1
			f /= 10
		}
		nums = append(nums, f)
	}
	return nums
}

func GetSumOfAmounts(nums []float64) float64 {
	sum := float64(0)
	for i := 0; i < len(nums); i++ {
		sum = RoundToTwoDecimalPlaces(sum + RoundToTwoDecimalPlaces(nums[i]))
	}
	return sum
}
