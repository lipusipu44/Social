package main

import (
	"context"
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
	Title   string   `json:"title" validate:"required,max=255"` //for validation purpose details below in methods
	Content string   `json:"content" validate:"required,max=255"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createPostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(res http.ResponseWriter, req *http.Request) {
	var payload createPostPayload
	if err := readJSON(res, req, &payload); err != nil {
		app.badRequestResponse(res, req, err)
		return
	}
	/*
		this block is actually responsible for validation check not the json section in struct
		this checks using that validate written in struct block of createPostPayload.
	*/
	if err := CustomValidate.Struct(&payload); err != nil {
		app.badRequestResponse(res, req, err)
		return
	}
	usr := getUserfromMiddleWare(req)
	/*
		It needs a min payload of Post struct which are required in Create method
		of post in insertion query.
	*/
	post := &store.Post{
		UserID:  usr.ID, //changed from hardcoded value to get the value from middleware
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
	if err := writeJSONWrapper(res, http.StatusOK, post); err != nil {
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

Update : In this branch it also gets the comment used on that post_id
*/

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		string	true	"Post ID"
//	@Success		200		{object}	store.Post
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postId} [get]
func (app *application) getPostById(res http.ResponseWriter, req *http.Request) {
	/*
		the logic is moved to the middleware section named postContextMiddleware which tops up
		the handler with post; below that there is a method getPostFromContext which extracts
		the Post struct from the handler by using req context, details explained in those 2 methods
	*/
	post := getPostFromContext(req)

	/*
		in this branch Post struct has got comment as a field,
		there we are storing comments in array for that post id if any
	*/
	comments, err := app.store.Comment.GetCommentOfUserOnPost(req.Context(), post.ID)
	if err != nil {
		app.internalServerError(res, req, err)
		return
	}

	post.Comment = comments
	if err := writeJSONWrapper(res, http.StatusOK, post); err != nil {
		app.internalServerError(res, req, err)
	}
}

//deletePostHandler
/*
simple delete post handler based on post id
*/
func (app *application) deletePostHandler(res http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParam(req, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestResponse(res, req, err)
		return
	}
	if err := app.store.Post.Delete(req.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrNoRows):
			app.notFoundResponse(res, req, err)
			return
		default:
			app.internalServerError(res, req, err)
			return
		}
	}
	/*
		no content to show as result, but status to be shown as 204 no content as all deleted
	*/
	res.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   string `json:"title" validate:"omitempty,max=255"`
	Content string `json:"content" validate:"omitempty,max=255"`
}

//updatePostHandler
/*
check the updated comment in Update Post in Post.go in store package,
there I have mentioned about concurrent user case, how they will block it byu comparing
version with DB version
*/
// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postId} [patch]
func (app *application) updatePostHandler(res http.ResponseWriter, req *http.Request) {
	post := getPostFromContext(req)
	var payload UpdatePostPayload
	if err := readJSON(res, req, &payload); err != nil {
		app.badRequestResponse(res, req, err)
	}
	if err := CustomValidate.Struct(&payload); err != nil {
		app.badRequestResponse(res, req, err)
	}
	/*
		if title or content are "" then it will update
		with original title or content or
		take the title and content and change it
	*/
	if payload.Title != "" {
		post.Title = payload.Title
	}
	if payload.Content != "" {
		post.Content = payload.Content
	}
	err, post := app.store.Post.Update(req.Context(), post) //post on Update we got from middleware
	if err != nil {
		switch {
		case errors.Is(err, store.ErrSerializationFailure):
			app.retryResponse(res, req, err) // You might want to implement a retry mechanism

			//all the below cases I did for dummy work, not tested so dont panic
		case errors.Is(err, store.ErrUniqueViolation):
			app.badRequestResponse(res, req, err)

		case errors.Is(err, store.ErrDeadlock):
			app.conflictResponse(res, req, err)

		default:
			app.internalServerError(res, req, err)

		}

	}
	if err := writeJSONWrapper(res, http.StatusOK, post); err != nil {
		app.internalServerError(res, req, err)
	}
}

/*
Middleware section starts from here, please read it 2 times for better understanding
*/
// Context key type to avoid conflicts, to be used in extract key from handler
type contextKey string

const postKey contextKey = "post"

//postContextMiddleware
/*
In Go (Golang), a handler function is used in web development to handle HTTP requests

typically middleware is a medium by which handler tops-itself up with
new information, in this case its Post struct if the logic is able to find a post-struct, then middleware
handler take it inside itself.

next handler is kind of boggy which now will contain Post struct in itself, so wherever
another handler call will happen, it will take the data out of this handler and use it.
*/
func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
		//above section direct copy paste from GET-Post call

		/*
			below line context adds post struct with a key postKey,
			this key will be used in below method to extract Post
			struct from req context.
		*/
		ctx := context.WithValue(req.Context(), postKey, post)

		// Pass modified request to the handler named next
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

//getPostFromContext
/*
by using above postContextMiddleware next handler will push the data to r http.Request,
from there by using contextKey we are extracting the data/ Post struct in this method
*/
func getPostFromContext(r *http.Request) *store.Post {
	postExtract, _ := r.Context().Value(postKey).(*store.Post)
	return postExtract
}
