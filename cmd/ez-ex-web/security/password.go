package security

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// HashedPassword provides information on the hash of a password
type HashedPassword struct {
	Alg  string
	Hash []byte
	Salt []byte
	Ver  uint32
	Mem  uint32
	Time uint32
	CPUs uint8
}

type PasswordHasher interface {
	WithConfig(time, mem uint32, cpus uint8) PasswordSalter
	Decode() (*HashedPassword, error)
}

type PasswordSalter interface {
	WithSaltSize(size uint32) PasswordGenerator
}

type PasswordGenerator interface {
	Hash() HashEncoder
	Verify(hash, salt []byte) bool
}

type HashEncoder interface {
	Encode() string
	Bytes() (hash []byte, salt []byte)
}

type passwordHasher struct {
	pwd  []byte
	salt struct {
		size uint32
		b    []byte
	}
	hash struct {
		time uint32
		mem  uint32
		cpus uint8
		b    []byte
	}
}

// NewPasswordHasher initializes a new Argon2id password hasher
func NewPasswordHasher(pwd string) PasswordHasher {
	return &passwordHasher{
		pwd: []byte(pwd),
	}
}

// WithConfig sets Argon2id time (number of iterations) and memory usage
func (ph *passwordHasher) WithConfig(time, mem uint32, cpus uint8) PasswordSalter {
	ph.hash.time = time
	ph.hash.mem = mem
	ph.hash.cpus = cpus
	return ph
}

// WithSaltSize sets the salt size used during hashing
func (ph *passwordHasher) WithSaltSize(size uint32) PasswordGenerator {
	ph.salt.size = size
	return ph
}

// Hash hashes and salts the password
func (ph *passwordHasher) Hash() HashEncoder {
	ph.hash.b, ph.salt.b = HashPassword(ph.pwd, ph.salt.size, ph.hash.time, ph.hash.mem, ph.hash.cpus)
	return ph
}

// Bytes returns hash and salt bytes
func (ph *passwordHasher) Bytes() (hash []byte, salt []byte) {
	return ph.hash.b, ph.salt.b
}

// Verify checks if the hash is valid
func (ph *passwordHasher) Verify(hash, salt []byte) bool {
	return VerifyHash(ph.pwd, hash, salt, ph.hash.time, ph.hash.mem, ph.hash.cpus)
}

// Encode encodes the resulting hash into a string
func (ph *passwordHasher) Encode() string {
	return fmt.Sprintf(
		"$argon2i$v=19$m=%d,t=%d,p=%d$%s$%s",
		ph.hash.mem,
		ph.hash.time,
		ph.hash.cpus,
		base64.RawStdEncoding.EncodeToString(ph.salt.b),
		base64.RawStdEncoding.EncodeToString(ph.hash.b),
	)
}

// Decode decodes the hashed password in its components
func (ph *passwordHasher) Decode() (*HashedPassword, error) {
	if len(ph.pwd) == 0 {
		return nil, ErrInvalidEncodingFormat
	}

	chunks := strings.Split(string(ph.pwd[1:]), "$")
	if len(chunks) != 5 {
		return nil, ErrInvalidEncodingFormat
	}

	sub := strings.Split(chunks[2], ",")
	if len(sub) != 3 {
		return nil, ErrInvalidEncodingFormat
	}

	var ver, mem, time, cpu int
	var err error

	if len(chunks[1]) < 4 ||
		len(sub[0]) < 3 ||
		len(sub[1]) < 3 ||
		len(sub[2]) < 3 {
		return nil, ErrInvalidEncodingFormat
	}

	if ver, err = strconv.Atoi(chunks[1][2:]); err != nil {
		return nil, err
	}
	if mem, err = strconv.Atoi(sub[0][2:]); err != nil {
		return nil, err
	}
	if time, err = strconv.Atoi(sub[1][2:]); err != nil {
		return nil, err
	}
	if cpu, err = strconv.Atoi(sub[2][2:]); err != nil {
		return nil, err
	}

	hash := make([]byte, base64.RawStdEncoding.DecodedLen(len(chunks[4])))
	salt := make([]byte, base64.RawStdEncoding.DecodedLen(len(chunks[3])))
	if _, err = base64.RawStdEncoding.Decode(salt, []byte(chunks[3])); err != nil {
		return nil, err
	}
	if _, err = base64.RawStdEncoding.Decode(hash, []byte(chunks[4])); err != nil {
		return nil, err
	}

	return &HashedPassword{
		Alg:  chunks[0],
		Hash: hash,
		Salt: salt,
		Ver:  uint32(ver),
		Mem:  uint32(mem),
		Time: uint32(time),
		CPUs: uint8(cpu),
	}, err
}
