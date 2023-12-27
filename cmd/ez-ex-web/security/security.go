package security

import (
	"bytes"
	cryptoRand "crypto/rand"
	"golang.org/x/crypto/argon2"
)

func generateSalt(s uint32) []byte {
	var salt = make([]byte, s)
	_, _ = cryptoRand.Read(salt)

	return salt
}

func hashAndSalt(pwd, salt []byte, hashTime, hashMem uint32, cpus uint8) []byte {
	return argon2.IDKey(pwd, salt, hashTime, hashMem, cpus, uint32(len(salt)))
}

func hashPassword(pwd []byte, saltSize, hashTime uint32, hashMem uint32, cpus uint8) (hash []byte, salt []byte) {
	salt = generateSalt(saltSize)
	hash = hashAndSalt(pwd, salt, hashTime, hashMem, cpus)
	return hash, salt
}

func verifyHash(pwd, hash, salt []byte, hashTime uint32, hashMem uint32, cpus uint8) bool {
	if salt == nil || len(salt) == 0 {
		return false
	}

	res := hashAndSalt(pwd, salt, hashTime, hashMem, cpus)
	return bytes.Equal(res, hash)
}
