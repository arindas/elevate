package main

import "flag"

var baseDirectory = flag.String("base_dir", ".", "Base directory for storing files.")

func main() {
	flag.Parse()

	ServeApp(AppConfig{BaseDirectory: *baseDirectory})
}
