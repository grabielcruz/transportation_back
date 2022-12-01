package utility

import "testing"

func TestRandomString(t *testing.T) {
	randStr1 := GetRandomString(23)
	randStr2 := GetRandomString(23)
	if randStr1 == randStr2 {
		t.Fatalf(`GetRandomString() generated the same random string twice in a row, randStr1 = %v, randStr2 = %v`, randStr1, randStr2)
	}
}

func TestRandomBalance(t *testing.T) {
	randBal1 := GetRandomBalance()
	randBal2 := GetRandomBalance()
	if randBal1 == randBal2 {
		t.Fatalf(`GetRandomBalance() generated the same random balance twice in a row, randBal1 = %v, randBal2 = %v`, randBal1, randBal2)
	}
}

func TestGetRandomBoolean(t *testing.T) {
	passed := false
	for i := 0; i < 100; i++ {
		randBool1 := GetRandomBoolean()
		randBool2 := GetRandomBoolean()
		if randBool1 != randBool2 {
			passed = true
			break
		}
	}
	if !passed {
		t.Fatalf(`GetRandomBoolean() dit not generate different booleans in 100 iterations`)
	}
}

func TestGetRandomCurrency(t *testing.T) {
	passed := false
	for i := 0; i < 100; i++ {
		randCurr1 := GetRandomCurrency()
		randCurr2 := GetRandomCurrency()
		if randCurr1 != randCurr2 {
			passed = true
			break
		}
	}
	if !passed {
		t.Fatalf(`TestGetRandomCurrency() dit not generate different currencies in 100 iterations`)
	}
}
