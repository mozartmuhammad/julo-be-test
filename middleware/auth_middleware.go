package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/mozartmuhammad/julo-be-test/helper"
)

// Define JWTClaimsContextKey as the key to access JWT claims in the context
const JWTClaimsContextKey = "jwtClaims"
const CustomerXIDContextKey = "customer_xid"

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	CustomerXID string `json:"customer_xid"`
	jwt.StandardClaims
}

func AuthorizeRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Read secret key from environment variable
		secretKey := os.Getenv("SECRET")
		if secretKey == "" {
			fmt.Println("SECRET is not set in .env file")
		}

		tokenString := strings.Split(authHeader, " ")[1]
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := helper.SetCustomerXID(r.Context(), claims.CustomerXID)
		fn(w, r.WithContext(ctx))
	}
}
