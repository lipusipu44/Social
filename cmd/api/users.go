package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lipusipu44/Social/internal/store"
	"log"
	"net/http"
	"strconv"
)

//getUserById goDoc

// @Summary		Get user by ID
// @Description	Fetches the user details from the middleware and returns the user data in JSON format.
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			userId	path		string		true	"User ID"
// @Success		200		{object}	store.User	"Successfully retrieved user"
// @Failure		400		{object}	error		"Bad request"
// @Failure		404		{object}	error		"User not found"
// @Failure		500		{object}	error		"Internal server error"
// @Security		ApiKeyAuth
// @Router			/users/{userId} [get]
func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	userStruct := getUserfromMiddleWare(r)

	if err := writeJSONWrapper(w, http.StatusOK, userStruct); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

}

type Follower struct {
	UserId int64 `json:"user_id"`
}

/*
logic self-explanatory, above struct temporarily used, once auth
is in place these might not be required
*/
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromMiddleWare(r)
	var followerUser Follower
	err := readJSON(w, r, &followerUser)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	errUnique := app.store.Follow.Follow(r.Context(), user.ID, followerUser.UserId)
	if errUnique != nil {
		switch {
		case errors.Is(err, store.ErrUniqueViolation):
			//unique key check and working successfully
			app.uniqueContraintConflictResponse(w, r, errUnique)
			return
		default:
			app.internalServerError(w, r, errUnique)
			return
		}
	}
	if err := writeJSONWrapper(w, http.StatusNoContent, user); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			tokenId	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{tokenId} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "tokenId")
	log.Println("Token: ", token)
	if err := app.store.User.Activate(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, store.ErrNoRows):
			app.badRequestResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := writeJSONWrapper(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

/*
Below lines are for middleware.
same concept as we applied for POST,
here is almost the same concept for middleware of user
*/
type userContextKey string

var userKey userContextKey = "user"

func (app *application) userMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "userId")
		ctx := r.Context()
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		userStruct, err := app.store.User.GetByID(r.Context(), userId)
		if err != nil {
			switch err {
			case store.ErrNoRows:
				app.notFoundResponse(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, userKey, userStruct)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserfromMiddleWare(r *http.Request) *store.User {
	userExtract, _ := r.Context().Value(userKey).(*store.User)
	return userExtract
}
