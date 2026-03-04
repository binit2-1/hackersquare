package server

import (
	"context"
	"net/http"
	"os"

	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
)



func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized: No auth token provided", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(cookie.Value, os.Getenv("JWT_SECRET"))
		if err != nil {
			http.Error(w, "Unauthorized: Invalid auth token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID",claims.UserID)
		reqWithContext := r.WithContext(ctx)

		next.ServeHTTP(w, reqWithContext)
	}
}