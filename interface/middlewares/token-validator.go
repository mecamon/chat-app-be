package middlewares

import (
	"context"
	json_web_token "github.com/mecamon/chat-app-be/interface/json-web-token"
	"net/http"
)

func TokenValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customClaims, err := json_web_token.Validate(r.Header.Get("Authorization"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(nil)
			return
		}

		ctx := context.WithValue(r.Context(), "ID", customClaims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
