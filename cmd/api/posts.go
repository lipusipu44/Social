package main

import (
	"github.com/lipusipu44/Social/internal/store"
	"net/http"
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
		writeJSONError(res, http.StatusBadRequest, err.Error())
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
		In the below statement we are sending the partial post,
		but it will get the rest of the value from create scan method and
		store it in post, that's why reference is passed, and we will send it
		to writeJSON for full Post payload creation as response json
	*/
	/*
		this piece is critical and always confusing for me
		app.store.Post.Create
		app meaning object of application which came from method app *application
		app struct has a key called store - this is assigned with PostStorage and UserStorage
		by line - store.NewStorage(db) in main.go,
		store struct has a key called Post, its not Post class or anything :) which has value as Post Interface,
		there always confusion
		instead of using Post interface its using PostStorage as its implements Post interface. in PostStorage.go
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
		writeJSONError(res, http.StatusInternalServerError, err.Error())
		return
	}
	if err := writeJSON(res, http.StatusOK, post); err != nil {
		writeJSONError(res, http.StatusInternalServerError, err.Error())
		return
	}
}
