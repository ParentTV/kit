package encoding

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const testSecret = "test-secret-not-used-anywhere"

func sign(t *testing.T, method jwt.SigningMethod, key interface{}, claims jwt.MapClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{"sub": "user-1", "exp": float64(time.Now().Add(time.Hour).Unix())}
}

func TestParseTokenString(t *testing.T) {
	t.Run("accepts a token signed with JWT_SECRET", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testSecret)
		a := ParseTokenString(sign(t, jwt.SigningMethodHS256, []byte(testSecret), validClaims()))
		if a.Error != nil || a.User != "user-1" || a.Expiry.IsZero() {
			t.Fatalf("got %+v", a)
		}
	})

	t.Run("rejects a token signed with another secret", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testSecret)
		a := ParseTokenString(sign(t, jwt.SigningMethodHS256, []byte("someone-elses"), validClaims()))
		if a.Error == nil || a.User != "" {
			t.Fatalf("got %+v", a)
		}
	})

	t.Run("rejects everything when JWT_SECRET is unset, even empty-key tokens", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "")
		for _, key := range [][]byte{[]byte(testSecret), {}} {
			a := ParseTokenString(sign(t, jwt.SigningMethodHS256, key, validClaims()))
			if !errors.Is(a.Error, ErrNoJWTSecret) || a.User != "" {
				t.Fatalf("got %+v", a)
			}
		}
	})

	t.Run("rejects alg none", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testSecret)
		a := ParseTokenString(sign(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, validClaims()))
		if a.Error == nil || a.User != "" {
			t.Fatalf("got %+v", a)
		}
	})

	t.Run("rejects expired tokens", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testSecret)
		c := validClaims()
		c["exp"] = float64(time.Now().Add(-time.Minute).Unix())
		if a := ParseTokenString(sign(t, jwt.SigningMethodHS256, []byte(testSecret), c)); a.Error == nil {
			t.Fatalf("got %+v", a)
		}
	})

	t.Run("a signed token without a subject is an error, not a panic", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testSecret)
		a := ParseTokenString(sign(t, jwt.SigningMethodHS256, []byte(testSecret), jwt.MapClaims{"exp": float64(time.Now().Add(time.Hour).Unix())}))
		if a.Error == nil || a.User != "" {
			t.Fatalf("got %+v", a)
		}
	})
}

func TestParseToken(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+sign(t, jwt.SigningMethodHS256, []byte(testSecret), validClaims()))
	if a := ParseToken(r); a.User != "user-1" {
		t.Fatalf("got %+v", a)
	}
	if a := ParseToken(httptest.NewRequest("GET", "/", nil)); a.User != "" || a.Error != nil {
		t.Fatalf("no header should be an empty Auth, got %+v", a)
	}
}
