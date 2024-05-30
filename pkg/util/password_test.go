package util_test

import (
	"strings"
	"testing"

	"github.com/8thgencore/passfort/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("Successfully hashes password", func(t *testing.T) {
		password := "test-password"

		hashedPassword, err := util.HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("Returns error when hashing fails", func(t *testing.T) {
		password := strings.Repeat("a", 100)

		_, err := util.HashPassword(password)
		assert.Error(t, err)
	})
}

func TestCompareHash(t *testing.T) {
	t.Run("Successfully compares password with hashed password", func(t *testing.T) {
		password := "test-password"
		hashedPassword, err := util.HashPassword(password)
		assert.NoError(t, err)

		err = util.CompareHash(password, hashedPassword)
		assert.NoError(t, err)
	})

	t.Run("Returns error when password does not match hashed password", func(t *testing.T) {
		password := "test-password"
		hashedPassword, err := util.HashPassword(password)
		assert.NoError(t, err)

		err = util.CompareHash("wrong-password", hashedPassword)
		assert.Error(t, err)
	})
}
