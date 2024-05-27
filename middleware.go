package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

func validateJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
}

func createJWT(account *Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      account.ID,
		"expires": time.Now().AddDate(0, 0, 14),
	})

	tokenString, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func makeAuthMiddleware(s *APIServer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			json_user, err := json_stringify(&Account{ID: 1, FirstName: "Kez", LastName: "Anwar", Number: 123321, Balance: 0})

			if err != nil {
				log.Fatal(err)
			}

			r.Header.Add("x-server-user", json_user)

			if r.Method == "POST" && r.RequestURI == "/account" {
				next.ServeHTTP(w, r)
				return
			}

			reqToken := r.Header.Get("x-auth-token")

			if len(reqToken) == 0 {
				WriteJSON(w, http.StatusForbidden, ApiError{Message: "invalid token"})
				return
			}

			token, err := validateJWT(reqToken)

			if err != nil {
				WriteJSON(w, http.StatusForbidden, ApiError{Message: "invalid token"})
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
