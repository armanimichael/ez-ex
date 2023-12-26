package security

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordHasher_Panic(t *testing.T) {
	cases := []struct {
		time     uint32
		mem      uint32
		cpus     uint8
		saltSize uint32
	}{
		{0, 1, 1, 1},
		{1, 1, 0, 1},
	}

	for _, c := range cases {
		test := fmt.Sprintf(
			"WithConfig(%d, %d, %d), WithSaltSize(%d)",
			c.time,
			c.mem,
			c.cpus,
			c.saltSize,
		)
		t.Run(test, func(t *testing.T) {
			assert.Panics(t, func() {
				NewPasswordHasher("test").
					WithConfig(c.time, c.mem, c.cpus).
					WithSaltSize(c.saltSize).
					Hash()
			})
		})
	}
}

func TestPasswordHasher_Bytes(t *testing.T) {
	cases := []struct {
		time     uint32
		mem      uint32
		cpus     uint8
		saltSize uint32
	}{
		{1, 1, 1, 1},
		{10, 1, 1, 10},
		{100, 1, 1, 100},
	}

	for _, c := range cases {
		test := fmt.Sprintf(
			"WithConfig(%d, %d, %d), WithSaltSize(%d)",
			c.time,
			c.mem,
			c.cpus,
			c.saltSize,
		)
		t.Run(test, func(t *testing.T) {
			hash, salt := NewPasswordHasher("test").
				WithConfig(c.time, c.mem, c.cpus).
				WithSaltSize(c.saltSize).
				Hash().
				Bytes()

			assert.NotEmpty(t, hash)
			assert.NotEmpty(t, salt)
		})
	}
}

func TestPasswordHasher_Encode_Decode(t *testing.T) {
	cases := []struct {
		pwd      string
		time     uint32
		mem      uint32
		cpus     uint8
		saltSize uint32
	}{
		{"a", 1, 1, 1, 1},
		{"b", 10, 1, 1, 10},
		{"c", 100, 1, 1, 100},
	}

	for _, c := range cases {
		test := fmt.Sprintf(
			"WithConfig(%d, %d, %d), WithSaltSize(%d)",
			c.time,
			c.mem,
			c.cpus,
			c.saltSize,
		)
		t.Run(test, func(t *testing.T) {
			encode := NewPasswordHasher(c.pwd).
				WithConfig(c.time, c.mem, c.cpus).
				WithSaltSize(c.saltSize).
				Hash().
				Encode()

			decode, err := NewPasswordHasher(encode).Decode()

			assert.Nil(t, err)
			assert.Equal(t, "argon2i", decode.Alg)
			assert.NotEmpty(t, decode.Hash)
			assert.NotEmpty(t, decode.Salt)
			assert.Equal(t, uint32(19), decode.Ver)
			assert.Equal(t, c.mem, decode.Mem)
			assert.Equal(t, c.time, decode.Time)
			assert.Equal(t, c.cpus, decode.CPUs)
		})
	}
}

func TestPasswordHasher_Decode_Err(t *testing.T) {
	cases := []string{
		"",
		"-",
		"$-$-$-$-$-",
		"$-$--$-,-,-$-$-",
		"$-$-$-_-,-,-$-$-",
		"$-$-$-_-,-,-$-$-",
		"$-$-_-_$-_-,-_-,-_-$-$-",
		"$-_-$ab12$-_-,-_-,-_-$-$-",
		"$-_-$ab12$ab1,-_-,-_-$-$-",
		"$-_-$ab12$ab1,ab1,-_-$-$-",
		"$-_-$ab12$ab1,ab1,ab1$-$-",
		"$-_-$ab12$ab1,ab1,ab1$AA$-",
		"$-_-$ab12$ab1,ab1,ab1$-$AA",
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			decode, err := NewPasswordHasher(c).Decode()

			assert.Nil(t, decode)
			assert.Error(t, err)
		})
	}
}

func TestPasswordHasher_Verify(t *testing.T) {
	cases := []struct {
		pwd      string
		time     uint32
		mem      uint32
		cpus     uint8
		saltSize uint32
	}{
		{"a", 1, 1, 1, 1},
		{"b", 10, 1, 1, 10},
		{"c", 100, 1, 1, 100},
	}

	for _, c := range cases {
		test := fmt.Sprintf(
			"WithConfig(%d, %d, %d), WithSaltSize(%d)",
			c.time,
			c.mem,
			c.cpus,
			c.saltSize,
		)
		t.Run(test, func(t *testing.T) {
			hash, salt := NewPasswordHasher(c.pwd).
				WithConfig(c.time, c.mem, c.cpus).
				WithSaltSize(c.saltSize).
				Hash().
				Bytes()

			isValid := NewPasswordHasher(c.pwd).
				WithConfig(c.time, c.mem, c.cpus).
				WithSaltSize(c.saltSize).
				Verify(hash, salt)

			assert.True(t, isValid)
		})
	}
}

func TestPasswordHasher_Verify_False(t *testing.T) {
	hash, salt := NewPasswordHasher("test").
		WithConfig(1, 1, 1).
		WithSaltSize(1).
		Hash().
		Bytes()

	isValid := NewPasswordHasher("another-value").
		WithConfig(1, 1, 1).
		WithSaltSize(1).
		Verify(hash, salt)

	assert.False(t, isValid)
}
