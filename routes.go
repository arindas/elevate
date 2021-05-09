package main

func Routes(app AppConfig) []Route {
	return []Route{
		{"/", nil},
		{"/upload", UploadFileHandler(app)},
	}
}
