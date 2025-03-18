package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
	"strconv"
)

func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
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

	if err := writeJSONWrapper(w, http.StatusOK, userStruct); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

}
