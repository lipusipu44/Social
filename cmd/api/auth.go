package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
	"time"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=32"`
	Password string `json:"password" validate:"required,max=32"`
	Email    string `json:"email" validate:"required,email,max=72"`
}

//UserWithToken
/*
created this struct to capture the user
with the token so that can be used to activate it
using value
*/
type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	//capturing the payload to struct
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//validation of payload using validator
	if err := CustomValidate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	usr := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		//Check - I am not adding password, that logic to be added below
	}

	/*
		here the logic of password hashing to be implemented
		password of store.User is changing from plain string to struct, check in store.Password
		section now

		usr.Password also is set in Hash(), check the implementation
	*/
	if err := usr.Password.Hash(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	/*
		creating token string using uuid, import it
		plainToken to be sent to user in email, hashtoken of the same
		to be stored in DB, once the user authenticates, we will compare
		plaintoken with hashtoken and allow the user if validated
	*/
	planToken := uuid.New().String()

	//store encrypted plainToken in DB
	hash := sha256.Sum256([]byte(planToken))
	hashToken := hex.EncodeToString(hash[:])

	//create the user and invite
	err := app.store.User.CreateAndInvite(ctx, usr, hashToken, app.config.mailConf.exp)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestResponse(w, r, err)
			return
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return

		}

	}

	//send mail - todo

	/*
		used to get the user with token in payload response for temp purpose
	*/
	userwithToken := UserWithToken{
		User:  usr,
		Token: planToken,
	}

	//if all good then write it to the response using writer handler
	if err := writeJSONWrapper(w, http.StatusCreated, userwithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=72"`
	Password string `json:"password" validate:"required,max=32"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	//parse payload credential
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := CustomValidate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//check if user exists with the given email
	user, err := app.store.User.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNoRows):
			app.unauthorizedErrorResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	//if user is there generate token and add claim
	/*
		After user is found we take its user id to create the claim using secrets create by our app,
		user.ID, expiry time. this claim to be used below to create the JWT token

		I checked the token details after decoding it and it looks perfect, sub was matching with user id
	*/
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.jwtConfiguration.expDate).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.jwtConfiguration.issuer,
		"aud": app.config.auth.jwtConfiguration.issuer,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	//if all good then write it to the response using writer handler
	if err := writeJSONWrapper(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}

}
