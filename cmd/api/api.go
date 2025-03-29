package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lipusipu44/Social/docs" //required to generate swagger doc
	"github.com/lipusipu44/Social/internal/store"
	"go.uber.org/zap"

	//imported the middleware for swagger
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

type mailConfig struct {
	exp time.Duration
}
type config struct {
	addr     string
	dbConfig dbConfig
	//as of now below I am using it in health check API response
	env string
	//added for swagger doc
	apiURL   string
	mailConf mailConfig // all details of mail
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
	//adding ZAP logger object here, so anyone in handler and cmd can use it
	zapLogger *zap.Logger
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

	// A good base middleware stack, no explanation, used GPT
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		/*
			Below section, we are doing a setup for swagger doc initiation
			all swag init doc to be there in the docs folder in framework

			the first description of swagger is added in main.go's main(), thats kinda
			intro for swagger doc

			gen-docs - this part is added in makefile to generate swagger doc for cmd and internal
			package in docs folder in this project.

			make gen-docs - created the doc in docs package, this is going to be a part of air.toml file to
			start it for each change for new doc creation

			pre_cmd = ["make gen-docs"] - to be added in air.toml file to create the docs everytime air
			starts running

			if all goes good - http://localhost:8080/v1/swagger/index.html will open basic swagger page
			I think this localhost:8080 comes from docs.SwaggerInfo.Host =apiURL
			which is mentioned in config, this is initiated in run()
			of api.go
		*/
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL), //The url pointing to API definition
		))
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postId}", func(r chi.Router) {
				//usage of handler details explained in posts.go file
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostById)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{tokenId}", app.activateUserHandler)
			r.Route("/{userId}", func(r chi.Router) {
				//use of get user from id middleware
				r.Use(app.userMiddleWare)
				r.Get("/", app.getUserById)
				r.Put("/follow", app.followUserHandler)
			})
			//read more about group
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})

		})

		//Future only public route, moved it out of authentication
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
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
	/*
		imp - first run swagger init, then the doc will be created in docs folder,
		then add docs.swaggerinfo section here and also import the doc folder
		in import section on top, ensure doc section is not by-default swagger given folder
		I changed it to github.com/lipusipu44/Social/docs
	*/
	//swag doc related
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

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
