package http

import "github.com/arindas/elevate/internal/app"

func Routes(app app.AppConfig) []Route {
	return []Route{
		{"/", nil},
		{"/upload", UploadFileHandler(app)},
	}
}
