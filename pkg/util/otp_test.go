package util_test

import (
	"testing"

	"github.com/8thgencore/passfort/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP(t *testing.T) {
	t.Run("Successfully generates OTP", func(t *testing.T) {
		otp := util.GenerateOTP()
		assert.Len(t, otp, 6)
		assert.True(t, util.ValidateOTP(otp))
	})
}

func TestValidateOTP(t *testing.T) {
	t.Run("Valid OTP", func(t *testing.T) {
		otp := "123456"
		assert.True(t, util.ValidateOTP(otp))
	})

	t.Run("Invalid OTP (length not 6)", func(t *testing.T) {
		otp := "12345"
		assert.False(t, util.ValidateOTP(otp))
	})

	t.Run("Invalid OTP (contains non-digit characters)", func(t *testing.T) {
		otp := "12345a"
		assert.False(t, util.ValidateOTP(otp))
	})
}
