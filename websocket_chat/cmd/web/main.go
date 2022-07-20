package main

import (
	"log"
	"net/http"

	"github.com/mbeaver502/LearningGolang_Sawler/websocket_chat/internal/handlers"
)

const SERVER_ADDRESS = ":8080"

func main() {
	mux := routes()

	log.Println("Starting WebSocket channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting server on", SERVER_ADDRESS)
	_ = http.ListenAndServe(SERVER_ADDRESS, mux)
}
