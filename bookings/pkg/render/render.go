package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/pkg/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var cache map[string]*template.Template

	// get template cache from app config -- create once, read many times!
	if app.UseCache {
		cache = app.TemplateCache
	} else {
		cache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := cache[tmpl]
	if !ok {
		log.Fatalln("failed to get template from cache")
	}

	// let's find out if there's an error when we execute the cached value
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

// createTemplateCache_v2 creates a template cache
// this will automatically add all our templates and layouts
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all files *.page.tmpl from ./templates/
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		log.Println(err)
		return cache, err
	}

	// range over all *.page.tmpl files
	for _, page := range pages {
		// get the name of the file itself and parse it as a template with that name
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println(err)
			return cache, err
		}

		// get all our layouts
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			log.Println(err)
			return cache, err
		}

		// associate any layouts with templates that require them
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Println(err)
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}

// AddDefaultData adds default data to the value pointed to by td.
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)

	return td
}
