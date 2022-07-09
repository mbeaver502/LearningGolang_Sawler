package main

import (
	"log"
	"net/http"

	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/config"
	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/handlers"
	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/render"
)

const PORT_NUMBER = ":8080"

// main is the program entrypoint.
func main() {
	app := setupAppConfig()

	repo := handlers.NewRepo(app)
	handlers.NewHandlers(repo)

	render.NewTemplates(app)

	// Launch a server and start serving requests on localhost:PORT_NUMBER
	log.Printf("Starting server on port %v\n", PORT_NUMBER)

	srv := &http.Server{
		Addr:    PORT_NUMBER,
		Handler: routes(app),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func setupAppConfig() *config.AppConfig {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache_v2()
	if err != nil {
		log.Fatalln(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	return &app
}
