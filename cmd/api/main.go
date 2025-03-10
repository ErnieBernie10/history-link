package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"historylink/internal/features/link"
	"historylink/internal/features/record"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	//port := os.Getenv("PORT")
	connStr := os.Getenv("DATABASE_URL")

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			// Create a new router & API
			router := http.NewServeMux()
			config := huma.DefaultConfig("history-link", "1.0.0")
			config.DocsPath = ""
			api := humago.New(router, config)

			conn, err := sql.Open("postgres", connStr)
			if err != nil {
				panic(err)
			}

			defer conn.Close()

			err = conn.Ping()
			if err != nil {
				panic(err)
			}

			router.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`))
			})

			// Configure the logger to include source information
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
			rs := record.NewRecordResources(conn, logger)
			ls := link.NewLinkResources(conn, logger)
			rs.MountRoutes(api)
			ls.MountRoutes(api)

			corsRouter := corsMiddleware(router)

			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), corsRouter)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
	log.Println("Shutdown...")
}
