package main

import "log"

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
	log.Fatal(app.run())
}
