package appconfig

import (
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dxps/user_dir_dgp/internal/handlers"
)

// Config defines the configuration of the application.
type Config struct {
	srv struct {
		port         int
		env          string
		readTimeout  string
		writeTimeout string
		idleTimeout  string
	}
	db struct {
		dsn string
	}
}

// App holds the configuration, database connection poll and the session manager instance.
type App struct {
	Config
	assetsFS   embed.FS
	httpServer *http.Server
	sessionMgr *scs.SessionManager
}

func InitApp(assetsFS embed.FS) *App {

	app := App{assetsFS: assetsFS}
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s - Help\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.IntVar(&app.srv.port, "srv-port", 9000, "HTTP listening port")
	flag.StringVar(&app.srv.env, "srv-env", "development", "Environment (dev|prod)")
	flag.StringVar(&app.srv.readTimeout, "srv-read-timeout", "30s", "HTTP request read-timeout")
	flag.StringVar(&app.srv.writeTimeout, "srv-write-timeout", "60s", "HTTP response write-timeout")
	flag.StringVar(&app.srv.idleTimeout, "srv-idle-timeout", "90s", "HTTP keep-alive idle-timeout")
	flag.StringVar(&app.db.dsn, "db-dsn", "", "Database datasource name")
	flag.Parse()

	sessionMgr := scs.New()
	sessionMgr.Lifetime = 24 * time.Hour
	sessionMgr.Cookie.Secure = false
	app.sessionMgr = sessionMgr

	app.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", app.srv.port),
		ReadTimeout:  app.parseStringDuration(app.srv.readTimeout),
		WriteTimeout: app.parseStringDuration(app.srv.writeTimeout),
		IdleTimeout:  app.parseStringDuration(app.srv.idleTimeout),
		Handler:      app.router(),
	}

	return &app
}

// parseStrindDuration will parse a string to a time.Duration format,
// setting to zero if the parsing fails.
func (app *App) parseStringDuration(d string) time.Duration {

	dur, err := time.ParseDuration(d)
	if err != nil {
		slog.Error(fmt.Sprintf("Invalid duration from '%s' value. Setting the duration to zero.", d))
		dur = 0
	}
	return dur
}

func (app *App) router() http.Handler {

	mux := http.NewServeMux()

	assetsHandler := http.FileServerFS(app.assetsFS)
	mux.Handle("GET /assets/", assetsHandler)

	hh := handlers.HttpHandlers{}
	mux.HandleFunc("GET /{$}", hh.HomePageHandler)
	mux.HandleFunc("GET /login", hh.LoginPageHandler)

	return app.sessionMgr.LoadAndSave(mux)
}

func (app *App) StartHttpServer() {

	slog.Info(fmt.Sprintf("Listening on port %d ...", app.srv.port))
	err := app.httpServer.ListenAndServe()
	if err != nil {
		slog.Error("HTTP Server error", "error", err)
	}
}
