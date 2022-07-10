package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom form struct, embeds url.Values object.
type Form struct {
	url.Values
	Errors errors
}

// New initialized a Form struct.
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: errors{},
	}
}

// Has checks if form field is in POST and not empty.
func (f *Form) Has(field string, r *http.Request) bool {
	return r.Form.Get(field) != ""
}
