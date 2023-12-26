package security

import (
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"strconv"
	"time"
)

// UserClaim represents the information sent together with auth token
type UserClaim struct {
	UserID      int
	Fingerprint string
}

const (
	UserFingerprintClamKey = "ctx"
)

// CreateAuthJWT generated a NON-signed auth JWT for a user
func CreateAuthJWT(ttlInSec int, claims UserClaim) (token jwt.Token, err error) {
	// exp, iat
	now := time.Now()
	ttl := time.Duration(ttlInSec) * time.Second
	expiry := now.Add(ttl)

	// Standardize number types into strings - verification type mismatch
	userID := strconv.FormatUint(uint64(claims.UserID), 10)

	token, err = jwt.NewBuilder().
		IssuedAt(now).
		Expiration(expiry).
		Subject(userID).
		Claim(UserFingerprintClamKey, claims.Fingerprint).
		Build()

	if err != nil {
		return nil, err
	}

	return token, nil
}

// SignJWTWithHS256 cryptographically sign a JWT token with a given secret (SHA-256)
func SignJWTWithHS256(token jwt.Token, secret []byte) (signed []byte, err error) {
	signed, err = jwt.Sign(token, jwt.WithKey(jwa.HS256, secret))
	if err != nil {
		return nil, err
	}

	return signed, nil
}

// ValidateAuthRS256Token validates JWT authenticity
// if no claims are passed, only their existence will be evaluated
func ValidateAuthRS256Token(token []byte, secret []byte) (claims UserClaim, err error) {
	var tok jwt.Token
	tok, err = jwt.Parse(token, jwt.WithKey(jwa.HS256, secret))
	if err != nil {
		return claims, err
	}

	// Check required
	if err = jwt.Validate(tok,
		jwt.WithRequiredClaim(jwt.SubjectKey),
		jwt.WithRequiredClaim(jwt.IssuedAtKey),
		jwt.WithRequiredClaim(jwt.ExpirationKey),
		jwt.WithRequiredClaim(UserFingerprintClamKey),
	); err != nil {
		return claims, err
	}

	err = decodeClaims(tok, &claims)
	return claims, err
}

func decodeClaims(token jwt.Token, claims *UserClaim) error {
	fingerprint, ok := token.Get(UserFingerprintClamKey)
	if !ok {
		return ErrInvalidJWTToken
	}

	userID, err := strconv.Atoi(token.Subject())
	if err != nil {
		return ErrInvalidJWTToken
	}

	claims.UserID = userID
	claims.Fingerprint = fingerprint.(string)

	return nil
}
