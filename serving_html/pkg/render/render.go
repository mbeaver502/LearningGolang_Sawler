package render

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	const layoutPath = "./templates/base.layout.tmpl"
	templatePath := strings.Join([]string{"./templates/", tmpl}, "")

	parsedTemplate, _ := template.ParseFiles(layoutPath, templatePath)
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Println("renderTemplate error:", err)
		return
	}
}
