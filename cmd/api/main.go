package main

import (
	"github.com/go-playground/validator/v10"
	db2 "github.com/lipusipu44/Social/internal/db"
	"github.com/lipusipu44/Social/internal/env"
	"github.com/lipusipu44/Social/internal/store"
	"log"
)

/*
This class initiates cfg, application and call run() of app
which starts the http server.

this class is mainly for initialization and running main func
*/
const version = "0.0.1" //only used in health check api as of now

//CustomValidate
/*
Below 2 lines used for validation of payload and var is initialized
at the beginning of the go run using init, same as static block
in java

this var is going to be used in all the handler classes which has payload
to validate
*/
var CustomValidate *validator.Validate

func init() {
	CustomValidate = validator.New(validator.WithRequiredStructEnabled())
}

func main() {
	cfg := config{
		/*
			now hardcoded value to be replaced by env values
			each time any value changes in .envrc do direnv allow
		*/
		addr: env.GetEnv("API_ADDR", ":8081"),
		/*
			Adding dbConfig by fetching value from property file
		*/
		dbConfig: dbConfig{
			addr: env.GetEnv("DB_ADDR", "postgresql://admin:adminpassword@localhost/"+
				"socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetEnvInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetEnvInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetEnv("API_ENV", "dev"),
	}
	/*
		creating db instance from db.go in internal/db package,
		this will fetch the value from config for New() method in db.go
		post this is created, this instance to be passed to storage var

		Imp to note here the config stays in dbConfig and config struct, but usage happens independently
	*/
	db, err := db2.New(cfg.dbConfig.addr, cfg.dbConfig.maxOpenConns, cfg.dbConfig.maxIdleConns, cfg.dbConfig.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("DB Connection Pool Established")
	//creating SQL storage and then pass it to api struct
	storage := store.NewStorage(db)

	app := application{
		config: cfg,
		store:  storage,
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
