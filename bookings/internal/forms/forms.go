package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds url.Values object.
type Form struct {
	url.Values
	Errors errors
}

// New initializes a Form struct.
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: errors{},
	}
}

// Valid returns true if there are no errors on the form; otherwise false.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks whether the given fields have non-blank values.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

// MinLength checks that a given field meets the required minimum length.
func (f *Form) MinLength(field string, length int) bool {
	x := f.Values.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

// Has checks if form field is in POST and not empty.
func (f *Form) Has(field string) bool {
	x := f.Values.Get(field)

	if x == "" {
		return false
	}

	return true
}

// IsEmail cecks for valid email address.
func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
		return false
	}

	return true
}
