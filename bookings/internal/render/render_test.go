package render

import (
	"net/http"
	"testing"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
)

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	var ww *mockWriter

	err = RenderTemplate(ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser", err)
	}

	err = RenderTemplate(ww, r, "doesnotexist.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist", err)
	}

	app.UseCache = true
	err = RenderTemplate(ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser", err)
	}

	app.UseCache = false
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	app.Session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/whatever", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = app.Session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
