package main

import (
	"github.com/lipusipu44/Social/internal/env"
	"github.com/lipusipu44/Social/internal/store"
	"log"
)

/*
This class initiates cfg, application and call run() of app
which starts the http server.

this class is mainly for initialization and running main func
*/
func main() {
	cfg := config{
		/*
			now hardcoded value to be replaced by env values
			each time any value changes in .envrc do direnv allow
		*/
		addr: env.GetEnv("API_ADDR", ":8081"),
	}

	//creating SQL store and then pass it to api struct
	store := store.NewStorage(nil)
	app := application{
		config: cfg,
		store:  store,
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
