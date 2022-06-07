package encoding

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

type Auth struct {
	User     string
	Password string
}

func ParseAuthHeader(r *http.Request) Auth {
	u, p, ok := r.BasicAuth()
	if !ok {
		return Auth{}
	}
	return Auth{u, p}
}
func ParseToken(r *http.Request) Auth {
	secretKey := "5f2b5cdbe5194f10b3241568fe4e2b24"
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("invalid token")
		return Auth{}
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	return Auth{User: claims["sub"].(string)}
}
