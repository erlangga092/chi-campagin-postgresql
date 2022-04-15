package middleware

import (
	"context"
	"fmt"
	"funding-app/app/auth"
	"funding-app/app/helper"
	"funding-app/app/key"
	"funding-app/app/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(h http.Handler, authService auth.Service, userService user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		userID := fmt.Sprintf("%s", claim["user_id"])

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		authCtx := context.WithValue(ctx, key.CtxAuthKey{}, user)

		// serve to next route
		h.ServeHTTP(w, r.WithContext(authCtx))
	})
}
