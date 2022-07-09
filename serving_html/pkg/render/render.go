package render

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var templateCache = make(map[string]*template.Template)

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func RenderTemplateOld(w http.ResponseWriter, tmpl string) {
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
func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check if we already have the rendered template in templateCache
	if _, ok := templateCache[t]; !ok {
		// need to create the template
		err = createTemplateCache(t)

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

func createTemplateCache(t string) error {
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
