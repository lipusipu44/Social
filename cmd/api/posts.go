package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
	"strconv"
)

//createPostPayload
/*
As of now we need only this parts of payload
in curl payload, rest of the things for Post
struct to be created as hardcoded value and couple of
them to be fetched in PostStore create method

ex:

{
    "title":"anilpatro044 Post",
    "content":"Hey Guys this is my anilpatro044 post",
    "tags":["anilpatro044"]

}
*/
type createPostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(res http.ResponseWriter, req *http.Request) {
	var payload createPostPayload
	if err := readJSON(res, req, &payload); err != nil {
		app.badRequestResponse(res, req, err)
		return
	}
	/*
		It needs a min payload of Post struct which are required in Create method
		of post in insertion query.
	*/
	post := &store.Post{
		UserID:  1,
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
	}
	ctx := req.Context()
	/*
		In the below statement we are sending the partial Post struct.
		However, it will get the rest of the value from scan method in Create method and
		store it in post, that's why reference is passed, and we will send it
		to writeJSON for full Post payload creation as response json
	*/
	/*
		this piece is critical and always confusing for me
		app.store.Post.Create
		app meaning object of application which came from method app *application
		app struct has a key called store - this is assigned with PostStorage and UserStorage
		by line - store.NewStorage(db) in main.go, which gets PostStorage and UserStorage structs
		store struct has a key called Post Interface, its not Post class or anything :) which has value as Post Interface,
		there is always confusion
		instead of using Post interface its using PostStorage as its implement Post interface. in PostStorage.go
		so it has got PostStorage
		(which implements interface Post by using all Post Interface method)
		, which has a method called create in PostStorage class
		PostStore struct has a db key, its assigned the value in main.go
		db, err := db2.New(cfg.dbConfig.addr, cfg.dbConfig.maxOpenConns, cfg.dbConfig.maxIdleConns, cfg.dbConfig.maxIdleTime)
		storage := store.NewStorage(db)
		this poststore is passed in Create method in PostStorage class
		that's the way db is accessible in that method.

	*/
	if err := app.store.Post.Create(ctx, post); err != nil {
		app.internalServerError(res, req, err)
		return
	}
	if err := writeJSON(res, http.StatusOK, post); err != nil {
		app.internalServerError(res, req, err)
		return
	}
}

//getPostById
/*
this method gets the Post based on PostID, if its not there it gets
no row found my created error, which is created in storage.go

if found, it returns a Post reference and its converted to JSON in handler method
in posts.go in cmd/api package.
*/
func (app *application) getPostById(res http.ResponseWriter, req *http.Request) {
	//it gets the param postId from req, not from chi, chi is just a method to get it
	idParam := chi.URLParam(req, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestResponse(res, req, err)
		return
	}
	post, err := app.store.Post.GetByID(req.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNoRows):
			app.notFoundResponse(res, req, err)
		default:
			app.internalServerError(res, req, err)
		}
		return
	}
	if err := writeJSON(res, http.StatusOK, post); err != nil {
		app.internalServerError(res, req, err)
	}
}
