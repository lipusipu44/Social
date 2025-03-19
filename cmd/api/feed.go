package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feed, err := app.store.Post.GetUserFeed(ctx, int64(12))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSONWrapper(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
