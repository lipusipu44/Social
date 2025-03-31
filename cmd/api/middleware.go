package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
	"strconv"
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

type userContextKey string

var userKey userContextKey = "user"

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedTokenResponse(w, r, fmt.Errorf("missing 'Authorization' header"))
			return
		}

		//split it and get token part
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			app.unauthorizedTokenResponse(w, r, fmt.Errorf("auth token is malformed"))
			return
		}

		/*logic to capture the token string from variable-parts and validate it using
		validate token mentioned in internal package in jwt.go jwt struct
		*/
		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedTokenResponse(w, r, err)
			return
		}

		/*
			reverse engg, get the claim from token and by using claim's sub which was mapped
			to user.id we can get the user id, to check if he is valid or not
		*/
		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedTokenResponse(w, r, err)
			return
		}

		/*
			now that we got the userId lets validate if the user id is valid and
			active or not, very simple use app.store.user.GetByUserId
		*/

		user, err := app.store.User.GetByID(r.Context(), userID)
		if err != nil {
			app.unauthorizedTokenResponse(w, r, err)
			return
		}

		//normal middleware context logic
		ctx := context.WithValue(r.Context(), userKey, user) //userKey comes from users.go class
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getUserfromMiddleWare(r *http.Request) *store.User {
	userExtract, _ := r.Context().Value(userKey).(*store.User)
	return userExtract
}
