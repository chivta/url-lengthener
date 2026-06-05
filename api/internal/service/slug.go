package service

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

const (
	slugAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	slugLength   = 8
)

var validSlugRe = regexp.MustCompile(`^[a-zA-Z0-9]{1,16}$`)

func generateSlug() string {
	b := make([]byte, slugLength)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(slugAlphabet))))
		b[i] = slugAlphabet[n.Int64()]
	}
	return string(b)
}

func isValidSlug(s string) bool {
	return validSlugRe.MatchString(s)
}
