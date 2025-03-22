package main

import (
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//default pagination values, if nothing passed in URL
	fq := store.Pagination{
		Limit:  10,
		Offset: 0,
		SortBy: "desc",
	}
	//to parse the actual value from url param and put it in fq
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//validate the pagination struct or fail it
	if err := CustomValidate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	//value is temporarily hard coded, to be removed in auth section
	//change is - passing pagination struct to GetUserFeed
	feed, err := app.store.Post.GetUserFeed(ctx, int64(12), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSONWrapper(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
