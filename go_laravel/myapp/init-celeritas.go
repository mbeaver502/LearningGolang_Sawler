package main

import (
	"log"
	"os"

	"github.com/mbeaver502/LearningGolang_Sawler/go_laravel/celeritas"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	// initialize celeritas
	c := &celeritas.Celeritas{}
	err = c.New(path)
	if err != nil {
		log.Fatalln(err)
	}

	c.AppName = "myapp"
	c.Debug = true

	return &application{
		App: c,
	}
}
