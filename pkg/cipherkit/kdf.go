package cipherkit

import (
	"golang.org/x/crypto/argon2"
)

func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}
