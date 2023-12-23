package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

// auth is a middleware that parses a JWT from the Authorization header
// and extracts the userId from it. It then adds the userId to the request
// context. There's no actual authentication happening here 🤦‍♂️ something to
// do later.
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Print("No Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"no Authorization header"}`))
			return
		}
		vals := strings.Split(authHeader, "Bearer ")
		if len(vals) != 2 {
			log.Print("Invalid Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid Authorization header"}`))
		}
		tokenStr := vals[1]

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// todo: use a config manager like viper
			signingKey := os.Getenv("SIGNING_KEY")
			return []byte(signingKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Print("Invalid token signature")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"invalid token signature"}`))
				return
			}
			log.Print("Error parsing token")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"error parsing token"}`))
			return
		}

		if !tkn.Valid {
			log.Print("Invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Println("Valid user:", claims.UserId)

		// todo: can we use claims.sub instead of claims.UserId?
		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
