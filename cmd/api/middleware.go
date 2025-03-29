package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("missing 'Authorization' header"))
				return
			}
			//split it and get base64 part
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid 'Authorization' header"))
				return
			}

			//decode base64 part
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid 'Authorization' header"))
				return
			}

			pairs := strings.SplitN(string(decoded), ":", 2)
			username := app.config.auth.basicConfig.username
			pwd := app.config.auth.basicConfig.password

			if len(pairs) != 2 || pairs[0] != username || pairs[1] != pwd {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid 'Authorization' header"))
				return
			}
			next.ServeHTTP(w, r)

		})
	}
}
