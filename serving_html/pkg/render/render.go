package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/config"
	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/models"
)

var templateCache = make(map[string]*template.Template)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func RenderTemplate_v1(w http.ResponseWriter, tmpl string) {
	const layoutPath = "./templates/base.layout.tmpl"
	templatePath := strings.Join([]string{"./templates/", tmpl}, "")

	parsedTemplate, _ := template.ParseFiles(layoutPath, templatePath)
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Println("renderTemplate error:", err)
		return
	}
}

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func RenderTemplate_v2(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check if we already have the rendered template in templateCache
	if _, ok := templateCache[t]; !ok {
		// need to create the template
		err = createTemplateCache_v1(t)

		if err != nil {
			log.Println(err)
		}
	} else {
		// we already have template cached
		log.Println("using cached template", t)
	}

	tmpl = templateCache[t]
	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache_v1(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	// parse the template
	tmpl, err := template.ParseFiles(templates...)

	if err != nil {
		return err
	}

	// add template to cache
	templateCache[t] = tmpl

	return nil
}

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func RenderTemplate_v3(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// create a template cache -- this creates the cache on every render -- not good!
	//cache, err := CreateTemplateCache_v2()

	var cache map[string]*template.Template

	// get template cache from app config -- create once, read many times!
	if app.UseCache {
		cache = app.TemplateCache
	} else {
		cache, _ = CreateTemplateCache_v2()
	}

	// get requested template from cache
	t, ok := cache[tmpl]
	if !ok {
		log.Fatalln("failed to get template from cache")
	}

	// let's find out if there's an error when we execute the cached value
	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
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
func CreateTemplateCache_v2() (map[string]*template.Template, error) {
	//cache := make(map[string]*template.Template)
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

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}
