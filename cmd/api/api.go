package main

import (
	"log"
	"net/http"
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
The entry point for server to run, here everything is loaded for
http server and it starts running, it to be called in main.go file
*/
func (app *application) run() error {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}
	log.Printf("listening on %s", app.config.addr)
	return srv.ListenAndServe()
}
