package main

import (
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=32"`
	Password string `json:"password" validate:"required,max=32"`
	Email    string `json:"email" validate:"required,email,max=72"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"User registered"
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

	//create the user and invite
	err := app.store.User.CreateAndInvite(ctx, usr, "token-123")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	//send mail - todo

	//if all good then write it to the response using writer handler
	if err := writeJSONWrapper(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
