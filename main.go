package main

import (
	"log"
	"net/http"
)

type api struct {
	addr string
}

func (s *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			switch r.URL.Path {
			case "/":
				_, err := w.Write([]byte("index page"))
				if err != nil {
					return
				}
				return
			case "/users":
				_, err := w.Write([]byte("users page"))
				if err != nil {
					return
				}
				return
			}
		}
	default:
		_, err := w.Write([]byte("405 method not allowed"))
		if err != nil {
			return
		}
		return
	}
}

func main() {
	api := &api{
		addr: ":8080",
	}
	/*

			instead of using http.ListenAndServe,
			here we are declaring a server, to have a better control and
			from there we are calling ListenAndServe

			http.ListenAndServe(addr, handler) internally calls http.server {}
		and below method of srv.ListenAndServe

	*/
	srv := &http.Server{
		Addr:    api.addr,
		Handler: api,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
