package utils

import (
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	s := rand.NewSource(time.Now().UnixNano())
    r = rand.New(s)
}

func RandomInt(min, max int) int {
	return min + r.Intn(max-min)
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letters[r.Intn(len(letters))]
	}
	return string(bytes)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return int64(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[r.Intn(n)]
}