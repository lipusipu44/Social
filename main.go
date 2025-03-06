package main

import (
	"log"
	"net/http"
)

/*
Since we divided the logic into 3 files we should run
it using - go run *.go
*/
func main() {
	api := &api{
		addr: ":8080",
	}
	/*
		In this case we are using NewServeMux, this is basically same
		as the handler function which was written previously, it also implements
		ServeHTTP, so mux can be used in Handler based on Inheritance.

		this also takes care of better GET and POST logic handling.
		instead of bulk code in previous ServeHTTP way now we are breaking it to multiple parts
	*/
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}
	/*
		Also check the usage of HandleFunc by clicking on it.
		Self-explanatory
	*/
	mux.HandleFunc("GET /users", api.getUserHandler)
	mux.HandleFunc("POST /users", api.createUserHandler)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
