package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got invalid form when should have gotten valid form")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	if form.Has("whatever") {
		t.Error("shows has field when it does not")
	}

	postData = url.Values{}
	postData.Add("a", "a")

	form = New(postData)
	if !form.Has("a") {
		t.Error("shows does not have field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("shows valid when MinLength on non-existent field")
	}

	postData = url.Values{}
	postData.Add("some_field", "some value")

	form = New(postData)
	form.MinLength("some_field", 1)
	if !form.Valid() {
		t.Error("shows invalid when MinLength succeeds")
	}

	form = New(postData)
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows valid when MinLength fails")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("shows valid email when field does not exist")
	}

	postData = url.Values{}
	postData.Add("email", "test")

	form = New(postData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("shows valid for invalid email")
	}

	postData = url.Values{}
	postData.Add("email", "john@example.com")

	form = New(postData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("shows invalid for valid email")
	}
}
