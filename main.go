package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"

	"github.com/arindas/elevate/internal/app"
	"github.com/arindas/elevate/internal/http"
)

var baseDirectory = flag.String("base_dir", ".", "Base directory for storing files.")

//go:embed web/*
var webDirectory embed.FS

func main() {
	flag.Parse()

	log.Printf("Base directory: %s", *baseDirectory)

	staticContent, err := fs.Sub(webDirectory, "web")

	if err != nil {
		log.Printf("Unable to embed static files.")
	}

	server := http.ServerInstance(
		http.RequestHandler(
			http.Routes(app.AppConfig{
				BaseDirectory: *baseDirectory,
				StaticContent: staticContent,
			}),
			http.LoggingMiddleware,
		),
	)

	server.ReadServeAddr()
	server.ListenAndServe()
	server.Watch()

	server.LogErrors()
}
