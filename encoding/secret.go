package encoding

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// ErrNoJWTSecret is returned for every token when JWT_SECRET is unset. There is
// deliberately no default: an empty HMAC key would accept tokens anyone can sign.
var ErrNoJWTSecret = errors.New("JWT_SECRET is not set")

// jwtSecret is the HS256 key shared with identity-service, which signs the tokens.
// Read per call so a service picks it up without any wiring.
func jwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

// tokenKey only hands out the key for HMAC-signed tokens.
func tokenKey(t *jwt.Token, secret string) (interface{}, error) {
	if secret == "" {
		return nil, ErrNoJWTSecret
	}
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
	}
	return []byte(secret), nil
}
