package middleware

import (
	"context"
	"fmt"
	rs "go-jwt-api/response"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("abcdefghijklmnopq")

type Claims struct {
	Email              string `json:"email"`
	DisplayName        string `json:"displayName"`
	jwt.StandardClaims `json:"_"`
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			rs.ResponseErr(w, http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token
		if len(splitted) != 2 {
			rs.ResponseErr(w, http.StatusForbidden)
			return
		}

		tokenPart := splitted[1]
		fmt.Println("=================Token part===================" + tokenPart)
		tk := &Claims{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token không hợp lệ hoặc đã hết hạn", http.StatusUnauthorized)
			return
		}

		if token.Valid {
			ctx := context.WithValue(r.Context(), "user", tk)
			next(w, r.WithContext(ctx))

		}
	}

}
