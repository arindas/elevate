package app

import "io/fs"

type AppConfig struct {
	StaticContent fs.FS
	BaseDirectory string
}
