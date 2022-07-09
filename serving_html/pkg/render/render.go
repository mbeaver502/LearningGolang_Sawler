package render

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	path := strings.Join([]string{"./templates/", tmpl}, "")

	parsedTemplate, _ := template.ParseFiles(path)
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Println("renderTemplate error:", err)
		return
	}
}
