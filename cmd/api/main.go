package main

import (
	"github.com/go-playground/validator/v10"
	auth2 "github.com/lipusipu44/Social/internal/auth"
	db2 "github.com/lipusipu44/Social/internal/db"
	"github.com/lipusipu44/Social/internal/env"
	"github.com/lipusipu44/Social/internal/store"
	"go.uber.org/zap"
	"log"
	"time"
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

/*
Below all comments are used for swagger doc
*/
//	@title			Gopher Social API
//	@description	This is a Gopher Social.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	/*
		Introduction of proper structured logging using ZAP
	*/
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync() // to flush out the buffer

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
		//added for swagger doc
		apiURL: env.GetEnv("EXTERNAL_URL", "localhost:8081"),
		//adding mail config for invite user logic in auth.go
		mailConf: mailConfig{
			exp: time.Hour * 24 * 3,
		},
		//added config for basic config
		auth: authConfig{
			basicConfig: basic{
				username: env.GetEnv("API_USERNAME", "admin1"),
				password: env.GetEnv("API_PASSWORD", "admin1"),
			},
			//added all the jwt token configs from env files
			jwtConfiguration: jwtConfig{
				secret:  env.GetEnv("JWT_SECRET", "secret"),
				expDate: time.Hour * 24 * 3, // 3 days
				issuer:  env.GetEnv("JWT_ISSUER", "gophersocial"),
			},
		},
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
	logger.Info("DB Connection Pool Established")
	//creating SQL storage and then pass it to api struct
	storage := store.NewStorage(db)
	/*
		Creation of JWTAuthenticator and use it in app like storage is created
	*/
	jwtAuthenticator := auth2.NewJWTAuthenticator(cfg.auth.jwtConfiguration.secret,
		cfg.auth.jwtConfiguration.issuer,
		cfg.auth.jwtConfiguration.issuer)

	app := application{
		config: cfg,
		store:  storage,
		//passing logger to app
		zapLogger:     logger,
		authenticator: jwtAuthenticator,
	}
	//mount is initialized to accommodate HTTP calls

	/*
		Every method are attached to app,
		as app is initialized with all the config and
		in some form or others its going to be used in those
		methods
	*/
	mux := app.mount()
	logger.Fatal("Issue in starting APP server", zap.Error(app.run(mux)))
}
