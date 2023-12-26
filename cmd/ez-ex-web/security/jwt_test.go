package security

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
)

const (
	mockTTLInSec   = 10
	mockClaimKey   = "claim"
	mockClaimValue = "test"
)

var mockSecret = []byte("secret")
var mockSignedToken []byte

var mockUserClaims = UserClaim{
	UserID: 123,
}

func init() {
	tok, _ := CreateAuthJWT(mockTTLInSec, mockUserClaims)
	mockSignedToken, _ = SignJWTWithHS256(tok, mockSecret)
}

func TestSignJWTWithHS256(t *testing.T) {
	tok, _ := jwt.NewBuilder().Build()

	s, err := SignJWTWithHS256(tok, mockSecret)

	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

func TestSignJWTWithHS256_Structure(t *testing.T) {
	tok, _ := jwt.NewBuilder().Build()

	s, _ := SignJWTWithHS256(tok, mockSecret)
	parts := strings.Split(string(s), ".")

	assert.Len(t, parts, 3, "JWT must be composed of three parts")
}

func TestSignJWTWithHS256_Claims(t *testing.T) {
	tok, _ := jwt.NewBuilder().Claim(mockClaimKey, mockClaimValue).Build()

	s, _ := SignJWTWithHS256(tok, mockSecret)
	parsed, err := jwt.Parse(s, jwt.WithKey(jwa.HS256, mockSecret))
	c, ok := parsed.Get(mockClaimKey)

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, mockClaimValue, c)
}

func TestCreateAuthJWT(t *testing.T) {
	_, err := CreateAuthJWT(10, mockUserClaims)
	assert.NoError(t, err)
}

func TestCreateAuthJWT_Expiration(t *testing.T) {
	tok, _ := CreateAuthJWT(mockTTLInSec, mockUserClaims)
	expectedExpiry := tok.IssuedAt().Add(mockTTLInSec * time.Second).Unix()

	assert.Equal(t, expectedExpiry, tok.Expiration().Unix())
}

func TestCreateAuthJWT_UserID(t *testing.T) {
	tok, _ := CreateAuthJWT(mockTTLInSec, mockUserClaims)
	tokUserID, _ := strconv.Atoi(tok.Subject())

	assert.Equal(t, mockUserClaims.UserID, tokUserID)
}

func TestValidateAuthRS256Token(t *testing.T) {
	_, err := ValidateAuthRS256Token(mockSignedToken, mockSecret)
	assert.NoError(t, err)
}

func TestValidateAuthRS256Token_InvalidSecret(t *testing.T) {
	_, err := ValidateAuthRS256Token(mockSignedToken, []byte("this should not work"))
	assert.Error(t, err)
}

func TestValidateAuthRS256Token_Expired(t *testing.T) {
	tok, _ := CreateAuthJWT(-1, mockUserClaims)
	signed, _ := SignJWTWithHS256(tok, mockSecret)
	_, err := ValidateAuthRS256Token(signed, mockSecret)

	assert.Error(t, err)
}

func TestValidateAuthRS256Token_DecodedClaims(t *testing.T) {
	tok, _ := CreateAuthJWT(10, mockUserClaims)
	signed, _ := SignJWTWithHS256(tok, mockSecret)
	claims, _ := ValidateAuthRS256Token(signed, mockSecret)

	assert.Equal(t, mockUserClaims, claims)
}
