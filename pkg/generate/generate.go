package generate

import (
	"crypto/sha256"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// ID generates a unique identifier
func ID() string {
	return gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyz", 22)
}

// IDWithPrefix generates a unique identifier with a prefix.
func IDWithPrefix(prefix string) string {
	return prefix + ID()
}

// Token generates a random token.
func Token() string {
	return gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789*$%+&", 64)
}

// Hash generates a hash from a string.
func Hash(s string) (hash string) {
	sha256 := sha256.New()
	_, err := sha256.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", sha256.Sum(nil))
}
