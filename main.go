package main

import (
	"embed"
	"encoding/gob"
	"log/slog"
	"os"

	"github.com/dxps/user_dir_dgp/internal/appconfig"
)

//go:embed assets
var assetsFS embed.FS

// Sessions defines a custom typed to be used by the Session Manager.
type Sessions map[string]any

func init() {

	// Register the Sessions type, so the SessionManager can access it.
	gob.Register(&Sessions{})

	// Set the default slog logger.
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
}

func main() {

	app := appconfig.InitApp(assetsFS)
	app.StartHttpServer()
}
