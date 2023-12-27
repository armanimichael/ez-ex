package service

import (
	"context"
	"github.com/armanimichael/ez-ex/cmd/ez-ex-web/security"
)

type AuthService struct {
	SaltSize   uint32
	HashTime   uint32
	HashMemory uint32
	CPUs       uint8
}

func NewAuthService(saltSize uint32, hashTime uint32, hashMemory uint32, cpus uint8) AuthService {
	return AuthService{
		SaltSize:   saltSize,
		HashTime:   hashTime,
		HashMemory: hashMemory,
		CPUs:       cpus,
	}
}

// HashAndSaltPassword hashes a password with a random salt the encoded result
func (s AuthService) HashAndSaltPassword(_ context.Context, password string) string {
	return security.
		NewPasswordHasher(password).
		WithConfig(s.HashTime, s.HashMemory, s.CPUs).
		WithSaltSize(s.SaltSize).
		Hash().
		Encode()
}
