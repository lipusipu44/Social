package main

import (
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
From run section we are removing NewServerMux and creating
mount method, this creates server mux, attach the route like
GET, POST in it and then it will be used in run method
*/
func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
	return mux
}

/*
The entry point for server to run, here everything is loaded for
http server and it starts running, it to be called in main.go file
after mount is loaded with server mux
*/
func (app *application) run(mux *http.ServeMux) error {

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
