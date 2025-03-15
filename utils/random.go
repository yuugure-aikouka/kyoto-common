package utils

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// RandomString generates a random string of some length
func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomBase64EncodedString(length int) string {
	bytes := make([]byte, length)
	r.Read(bytes)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
}
