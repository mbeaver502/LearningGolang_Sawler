package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestGold_GetPrices(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer([]byte(jsonToReturn))),
		}
	})

	g := Gold{
		Prices: nil,
		Client: client,
	}

	p, err := g.GetPrices()
	if err != nil {
		t.Error(err)
	}

	if p.Price != 1772.005 {
		t.Errorf("got %v; want %v", p.Price, 1772.005)
	}

}
