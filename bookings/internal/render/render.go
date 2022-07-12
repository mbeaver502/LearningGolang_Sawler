package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
)

const (
	TEMPLATES_DIRECTORY = "./templates"
	TEMPLATE_FILE       = "*.page.tmpl"
	LAYOUT_FILE         = "*.layout.tmpl"
)

var pathToTemplates string = TEMPLATES_DIRECTORY

var app *config.AppConfig

// NewRenderer sets the config for the render package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// Template renders an HTML template to the given http.ResponseWriter.
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
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
		return errors.New("can't get template from cache")
	}

	// let's find out if there's an error when we execute the cached value
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}

	// render the template by writing to the passed-in http.ResponseWriter
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// createTemplateCache creates a template cache
// this will automatically add all our templates and layouts
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	templateFiles := fmt.Sprintf("%s/%s", pathToTemplates, TEMPLATE_FILE) //strings.Join([]string{TEMPLATES_DIRECTORY, TEMPLATE_FILE}, "")
	layoutFiles := fmt.Sprintf("%s/%s", pathToTemplates, LAYOUT_FILE)     // strings.Join([]string{TEMPLATES_DIRECTORY, LAYOUT_FILE}, "")

	// get all our page template files
	pages, err := filepath.Glob(templateFiles)
	if err != nil {
		log.Println(err)
		return cache, err
	}

	// range over all page template files
	for _, page := range pages {
		// get the name of the file itself and parse it as a template with that name
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println(err)
			return cache, err
		}

		// get all our layout template files
		layouts, err := filepath.Glob(layoutFiles)
		if err != nil {
			log.Println(err)
			return cache, err
		}

		// associate any layouts with templates that require them
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(layoutFiles)
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
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	return td
}
