package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	password        = "password"
	saltSize uint32 = 64
	hashTime uint32 = 3
	hashMem  uint32 = 64 * 1024
	hashCPUs uint8  = 1
)

func TestHashPassword_NonEmptyHash(t *testing.T) {
	pwd := []byte(password)
	hash, _ := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)

	assert.NotEmpty(t, hash)
}

func TestHashPassword_NonEmptySalt(t *testing.T) {
	pwd := []byte(password)
	_, salt := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)

	assert.NotEmpty(t, salt)
}

func TestHashPassword_PasswordNotEqualSalt(t *testing.T) {
	pwd := []byte(password)
	_, salt := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)

	assert.NotEqual(t, pwd, salt)
}

func TestHashPassword_PasswordNotEqualHash(t *testing.T) {
	pwd := []byte(password)
	hash, _ := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)

	assert.NotEqual(t, pwd, hash)
}

func TestHashPassword_HashNotEqualSalt(t *testing.T) {
	pwd := []byte(password)
	hash, salt := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)

	assert.NotEqual(t, hash, salt)
}

func TestVerifyHash_OK(t *testing.T) {
	pwd := []byte(password)
	pwdToVerify := []byte(password)

	hash, salt := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)
	isValid := VerifyHash(pwdToVerify, hash, salt, hashTime, hashMem, hashCPUs)

	assert.True(t, isValid)
}

func TestVerifyHash_WrongPassword(t *testing.T) {
	pwd := []byte(password)
	pwdToVerify := []byte("abc")

	hash, salt := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)
	isValid := VerifyHash(pwdToVerify, hash, salt, hashTime, hashMem, hashCPUs)

	assert.False(t, isValid)
}

func TestVerifyHash_WrongPassword_Empty(t *testing.T) {
	pwd := []byte(password)
	pwdToVerify := []byte("")

	hash, salt := HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)
	isValid := VerifyHash(pwdToVerify, hash, salt, hashTime, hashMem, hashCPUs)

	assert.False(t, isValid)
}

func TestVerifyHash_EmptySaltAndHash(t *testing.T) {
	var salt, hash []byte
	pwd := []byte(password)

	_, _ = HashPassword(pwd, saltSize, hashTime, hashMem, hashCPUs)
	isValid := VerifyHash(pwd, hash, salt, hashTime, hashMem, hashCPUs)

	assert.False(t, isValid)
}
