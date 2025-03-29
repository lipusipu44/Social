package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
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
