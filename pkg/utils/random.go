package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "1234567890"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomNumbers generates a random numbers of length n
func RandomNumbers(n int) string {
	var sb strings.Builder
	k := len(numbers)
	for i := 0; i < n; i++ {
		c := numbers[r.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
