package encoding

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	User     string
	Password string
	Expiry   time.Time
	Error    error
}

func ParseAuthHeader(r *http.Request) Auth {
	u, p, ok := r.BasicAuth()
	if !ok {
		return Auth{}
	}
	return Auth{User: u, Password: p}
}

func ParseToken(r *http.Request) Auth {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return Auth{}
	}
	reqToken = splitToken[1]
	return ParseTokenString(reqToken)
}

func ParseTokenString(t string) Auth {
	secretKey := jwtSecret()
	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		return tokenKey(t, secretKey)
	})
	if err != nil || !token.Valid {
		fmt.Println("invalid token:", err)
		return Auth{Error: err}
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(string)
	exp, _ := claims["exp"].(float64)
	if sub == "" {
		err = errors.New("token has no subject")
		fmt.Println("invalid token:", err)
		return Auth{Error: err}
	}
	return Auth{User: sub, Expiry: time.Unix(int64(exp), 0)}
}
