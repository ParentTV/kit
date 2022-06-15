package encoding

import (
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
	secretKey := "5f2b5cdbe5194f10b3241568fe4e2b24"
	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("invalid token")
		return Auth{}
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))
	return Auth{User: claims["sub"].(string), Expiry: time.Unix(exp, 0)}
}
