package main

import (
	"log"
)

/*
This class initiates cfg, application and call run() of app
which starts the http server.

this class is mainly for initialization and running main func
*/
func main() {
	cfg := config{
		addr: ":8080",
	}
	app := application{
		config: cfg,
	}
	//mount is initialized to accommodate HTTP calls

	/*
		Every method are attached to app,
		as app is initialized with all the config and
		in some form or others its going to be used in those
		methods
	*/
	mux := app.mount()
	log.Fatal(app.run(mux))
}
