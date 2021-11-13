package http

import (
	"net/http"

	"github.com/arindas/elevate/pkg/app"
)

func Routes(app app.AppConfig) []Route {
	return []Route{
		{"/", http.FileServer(http.FS(app.StaticContent))},
		{"/upload", UploadFileHandler(app)},
	}
}
