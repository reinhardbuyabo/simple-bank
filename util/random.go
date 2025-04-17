package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func init() {
	// Initialize the random number generator with a seed
	// This is important to ensure that the random numbers generated are not the same every time
	rand.Seed(time.Now().UnixNano())
}

// Generates a random integer between min and max (inclusive).
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generates a random string of length n using the characters in the alphabet.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code.
func RandomCurrency() string {
	currencies := []string{"KSH", "USD", "EUR", "CAD", "GBP", "JPY", "AUD", "CHF", "CNY", "SEK", "NZD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
