package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lipusipu44/Social/internal/store"
	"log"
	"net/http"
	"time"
)

/*
This api.go class is mainly used for any API related queries
starting from setting up the Api / application struct with diff
configs to start the API server.

Values to this struct to be injected from main class
*/

//dbConfig
/*
This contains all the config related to DB and properties to be fetched
from env property file, this struct is going to be a part
of config struct and initialization will happen from main.go
*/
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type config struct {
	addr     string
	dbConfig dbConfig
	//as of now below I am using it in health check API response
	env string
}

//application
/*
application struct contains all struct
as and when its required with function when app *application
attached it gets the fields of this struct

as and when web app becomes bigger application also contains
more and more struct
*/
type application struct {
	config config

	/*store is added in the struct
	so that it can be passed to handler methods and from
	there it will pass context to store layer, by which
	method will get the payload and params from context of handler
	*/
	store store.Storage //meaning Storage struct from store package
}

/*
going to use chi instead of serverMux
for api grouping based on
usage, and ease of usage of middleware
and authentication
*/

/*
chi.NewRouter() returns *chi.Mux which implements ServeHTTP,
so do handler too, so instead of returning chi.Mux we
would return http.Handler
*/
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack, no explanation use GPT
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", app.getPostById)
				r.Delete("/", app.deletePostHandler)
			})
		})
	})

	return r
}

/*
The entry point for server to run, here everything is loaded for
http server and it starts running, it to be called in main.go file
after mount is loaded with server mux
*/

/*
Imp - Here also return type instead of using *http.ServeMux we are using
http.Handler as both implements ServeHTTP method
*/
func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
		//Timeout configs
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("listening on %s", app.config.addr)
	return srv.ListenAndServe()
}
