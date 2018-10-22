package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/lamg/regapi"

	"github.com/lamg/ldaputil"
	_ "github.com/lib/pq"
	"gitlab.com/manuel.diaz/sisgguda/data"
)

type Configuration struct {
	ServerAddress      string
	ServerReadTimeout  int64
	ServerWriteTimeout int64

	UserDatabase     string
	NameDatabase     string
	PasswordDatabase string

	AccessLogPath string
	ErrorsLogPath string

	AdAddress  string
	AdSuff     string
	AdBDN      string
	AdUser     string
	AdPassword string
}

type App struct {
	Config    *Configuration
	datamodel *data.DataModel
	alogger   *log.Logger
	elogger   *log.Logger
	Handler   http.Handler
	ldap      *ldaputil.Ldap
	crypto    *regapi.JWTCrypt
}

func NewApp(config_path string) (app *App, err error) {
	app = new(App)

	//Loading configuration
	if err = app.loadConfig(config_path); err != nil {
		log.Fatalln("Cannot get configuration from file:", err)
	}

	//access_log
	accesslog_file, err := os.OpenFile(app.Config.AccessLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open access_log file:", err)
	}
	app.alogger = log.New(accesslog_file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	//errors_log
	errorslog_file, err := os.OpenFile(app.Config.ErrorsLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open errors_log file:", err)
	}
	app.elogger = log.New(errorslog_file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

	//Connecting to databse
	app.datamodel = new(data.DataModel)
	Db, err := sql.Open("postgres", "user="+app.Config.UserDatabase+" dbname="+app.Config.NameDatabase+" password="+app.Config.PasswordDatabase+" sslmode=disable")
	app.datamodel.Db = Db
	if err != nil {
		log.Fatalln("Failed to connect to database", err)
	}

	//Init Ldap
	app.ldap = ldaputil.NewLdapWithAcc(
		app.Config.AdAddress,
		app.Config.AdSuff,
		app.Config.AdBDN,
		app.Config.AdUser,
		app.Config.AdPassword)

	//Init crypto
	app.crypto = regapi.NewJWTCrypt()

	//Init Handlers
	app.RegisterHandlers()
	return
}

func (app *App) loadConfig(config_path string) (err error) {
	//Opennig config file
	file, err := os.Open(config_path)
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	defer file.Close()

	app.Config = &Configuration{}

	//Parsing json file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(app.Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}

	return
}

func (app *App) Name() string {
	return "SISGGUDA"
}

func (app *App) Version() string {
	return "0.1"
}

func (app *App) WriteResponse(w http.ResponseWriter, status int, cors bool, output []byte) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if cors {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Write(output)
}

const (
	RolSuperAdmin = "SuperAdmin"
	RolAdmin      = "Admin"
	RolSpecialist = "Specialist"
)

func (app *App) AllRols() []string {
	return []string{RolSuperAdmin, RolAdmin, RolSpecialist}
}

func (app *App) ValidRol(rol string) (err error) {
	if !containString(app.AllRols(), rol) {
		err = errors.New("Invalid Rol")
	}
	return
}

func (app *App) ValidAdUser(username string) (err error) {
	_, err = app.ldap.FullRecordAcc(username)
	return
}

func (app *App) RegisterHandlers() {
	router := mux.NewRouter()

	//Authentication
	router.Handle("/auth", &LoginHandler{app: app}).Methods(http.MethodPost)

	//User Operation Handlers
	router.Handle("/user", MustAuth(app, &AddUserHandler{app: app})).Methods(http.MethodPost)
	router.Handle("/user/{id:[0-9]+}", MustAuth(app, &UpdateUserHandler{app: app})).Methods(http.MethodPut)
	router.Handle("/user/{id:[0-9]+}", MustAuth(app, &DeleteUserHandler{app: app})).Methods(http.MethodDelete)

	//Group Operation Handlers
	router.Handle("/group", MustAuth(app, &AddGroupHandler{app: app})).Methods(http.MethodPost)
	router.Handle("/group/{id:[0-9]+}", MustAuth(app, &UpdateGroupHandler{app: app})).Methods(http.MethodPut)
	router.Handle("/group/{id:[0-9]+}", MustAuth(app, &DeleteGroupHandler{app: app})).Methods(http.MethodDelete)

	//GSpecialist Operation Handlers
	router.Handle("/gspecialist", MustAuth(app, &AddGroupSpecialistHandler{app: app})).Methods(http.MethodPost)
	router.Handle("/gspecialist", MustAuth(app, &DeleteGroupSpecialistHandler{app: app})).Methods(http.MethodDelete)

	/*
		//User Operation Handlers
		app.Router.Handle("/user", MustAuth(app, &GetUserHandler{app: app})).Methods(http.MethodGet)
		app.Router.Handle("/users", MustAuth(app, &GetUsersHandler{app: app})).Methods(http.MethodGet)

		//Group Operation Handlers
		app.Router.Handle("/group", MustAuth(app, GetGroupHandler{app: app})).Methods(http.MethodGet)
		app.Router.Handle("/groups", MustAuth(app, GetGroupsHandler{app: app})).Methods(http.MethodGet)

		//GSpecialist Operation Handlers
		app.Router.Handle("/gspecialist", MustAuth(app, AddGSpeecialistHandler{app: app})).Methods(http.MethodPost)
		app.Router.Handle("/gspecialist", MustAuth(app, UpdateGSpeecialistHandler{app: app})).Methods(http.MethodPut)
		app.Router.Handle("/gspecialist", MustAuth(app, DeleteGSpeecialistHandler{app: app})).Methods(http.MethodDelete)
		app.Router.Handle("/gspecialist", MustAuth(app, GetGSpeecialistHandler{app: app})).Methods(http.MethodGet)
		app.Router.Handle("/gspecialists", MustAuth(app, GetGSpeecialistsHandler{app: app})).Methods(http.MethodGet)

		//GUser Operation Handlers
		app.Router.Handle("/guser", MustAuth(app, AddGUserHandler{app: app})).Methods(http.MethodPost)
		app.Router.Handle("/guser", MustAuth(app, UpdateGUserHandler{app: app})).Methods(http.MethodPut)
		app.Router.Handle("/guser", MustAuth(app, DeleteGUserHandler{app: app})).Methods(http.MethodDelete)
		app.Router.Handle("/guser", MustAuth(app, GetGUserHandler{app: app})).Methods(http.MethodGet)
		app.Router.Handle("/gusers", MustAuth(app, GetGUserHandler{app: app})).Methods(http.MethodGet)
	*/

	app.Handler = cors.AllowAll().Handler(&AccessLogHandler{app: app, next: router})
}
