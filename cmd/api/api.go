package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

type config struct {
	addr string
}
type application struct {
	config config
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
