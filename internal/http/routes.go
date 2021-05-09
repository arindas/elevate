package http

import (
	"net/http"

	"github.com/arindas/elevate/internal/app"
)

func Routes(app app.AppConfig) []Route {
	return []Route{
		{"/", http.FileServer(http.Dir("./web"))},
		{"/upload", UploadFileHandler(app)},
	}
}
