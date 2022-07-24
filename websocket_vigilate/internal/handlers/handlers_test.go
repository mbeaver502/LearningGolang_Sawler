package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var loginTests = []struct {
	name       string
	url        string
	method     string
	postedData url.Values
	want       int //e.g., 200
}{
	{
		name:   "login-screen-get",
		url:    "/",
		method: "GET",
		want:   http.StatusOK,
	},
	{
		name:   "login-screen-post",
		url:    "/",
		method: "POST",
		postedData: url.Values{
			"email":    {"me@here.com"},
			"password": {"whatever"},
		},
		want: http.StatusSeeOther,
	},
}

func TestLoginScreen(t *testing.T) {
	for _, test := range loginTests {
		if test.method == "GET" {
			req, _ := http.NewRequest(test.method, test.url, nil)

			ctx := getCtx(req)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(Repo.LoginScreen)
			handler.ServeHTTP(rr, req)
			got := rr.Result().StatusCode

			if got != test.want {
				t.Errorf("%s got %d; want %d", test.name, got, test.want)
			}
		} else {
			req, _ := http.NewRequest(test.method, test.url, strings.NewReader(test.postedData.Encode()))

			ctx := getCtx(req)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(Repo.Login)
			handler.ServeHTTP(rr, req)
			got := rr.Result().StatusCode

			if got != test.want {
				t.Errorf("%s got %d; want %d", test.name, got, test.want)
			}
		}
	}
}
