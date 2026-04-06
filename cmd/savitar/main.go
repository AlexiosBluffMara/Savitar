package main

import (
	"os"

	"github.com/alexiosbluffmara/savitar/internal/app"
)

var version = "dev"

func main() {
	os.Exit(app.Run(os.Stdout, os.Stderr, version, os.Args[1:]))
}
