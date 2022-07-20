package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/mbeaver502/LearningGolang_Sawler/websocket_chat/internal/handlers"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	return mux
}