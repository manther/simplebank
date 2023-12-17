package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnop"
var r = *rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < k; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
