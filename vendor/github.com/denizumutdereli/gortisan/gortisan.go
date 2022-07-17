package gortisan

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/denizumutdereli/gortisan/render"
	"github.com/denizumutdereli/gortisan/session"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0-alpha"

type Gortisan struct {
	Appname  string `json:"appname"`
	Debug    bool   `json:"debug"`
	Version  string `json:"version"`
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string `json:"rootPath"`
	Routes   *chi.Mux
	Render   *render.Render
	Session  *scs.SessionManager
	DB       Database
	JetViews *jet.Set
	config   config
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    DatabaseConfig
}

func (g *Gortisan) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "caches", "data", "public", "tmp", "logs", "middleware"},
	}
	err := g.Init(pathConfig)
	if err != nil {
		return err
	}

	err = g.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	//create loggers
	infoLog, errorLog := g.startLoggers()

	//connect to database
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := g.OpenDB(os.Getenv("DATABASE_TYPE"), g.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		g.DB = Database{
			DatabaseType: os.Getenv("DATABASE_TYPE"),
			Pool:         db,
		}
	}

	g.InfoLog = infoLog
	g.ErrorLog = errorLog
	g.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	g.Version = version
	g.RootPath = rootPath
	g.Routes = g.routes().(*chi.Mux)

	g.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: DatabaseConfig{
			database: os.Getenv("DATABASE_NAME"),
			dsn:      g.BuildDSN(),
		},
	}

	// create session
	sess := session.Session{
		CookieLifeTime: g.config.cookie.lifetime,
		CookiePersist:  g.config.cookie.persist,
		CookieName:     g.config.cookie.name,
		CookieDomain:   g.config.cookie.domain,
		SessionType:    g.config.sessionType,
	}

	g.Session = sess.InitSessions()

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	g.JetViews = views

	g.createRenderer()

	return nil
}

func (g *Gortisan) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		//create the folder if doesn't not exist
		err := g.CreateDirIfNotExists(root + "/" + path)

		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gortisan) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     g.ErrorLog,
		Handler:      g.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	//close the DB connection at the end
	defer g.DB.Pool.Close()

	g.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	g.ErrorLog.Fatal(err)
}

func (g *Gortisan) checkDotEnv(path string) error {
	err := g.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))

	if err != nil {
		return err
	}

	return nil
}

func (g *Gortisan) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (g *Gortisan) createRenderer() {
	myRenderer := render.Render{
		Renderer: g.config.renderer,
		RootPath: g.RootPath,
		Port:     g.config.port,
		JetViews: g.JetViews,
	}

	g.Render = &myRenderer

}

func (g *Gortisan) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=%s connect_timeout=%s",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
			os.Getenv("DATABASE_TIMEZONE"),
			os.Getenv("DATABASE_CONNECT_TIMEOUT"),
		)

		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	default:

	}

	return dsn
}
