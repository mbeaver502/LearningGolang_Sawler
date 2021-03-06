package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

// views is a collection of jet views.
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

// upgradeConnection is a WebSocket connection upgrader provided by gorilla/websocket.
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Home renders the homepage.
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
		return
	}
}

// WebSocketConnection is a wrapper for the gorilla/websocket connection.
type WebSocketConnection struct {
	*websocket.Conn
}

// WsJSONResponse defines the JSON response from a WebSocket connection.
type WsJSONResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// WsPayload defines the JSON payload received over the WebSocket connection.
type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndpoint upgrades connection to WebSocket.
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade our HTTP connection to WebSocket
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("client connected to endpoint")

	var response WsJSONResponse
	response.Message = `<em><small>Connected to server</small></em>`

	// create a new wrapper for our newly upgraded WebSocket connection
	conn := WebSocketConnection{
		Conn: ws,
	}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
		return
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// no payload
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJSONResponse

	for {
		e := <-wsChan

		switch e.Action {

		// get a list of all users and send it back via broadcast
		case "username":
			// update the clients map to have the current user's username
			clients[e.Conn] = e.Username

			users := getUserList()

			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)

		// current user is leaving connection, so delete them and tell other users
		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			response.ConnectedUsers = getUserList()

			broadcastToAll(response)

		// handle user-sent messages
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)

			broadcastToAll(response)
		}

		// response.Action = "Got Here"
		// response.Message = fmt.Sprintf("Some Message, and action was %s", e.Action)
		// broadcastToAll(response)
	}
}

func getUserList() []string {
	var users []string

	for _, client := range clients {
		if client != "" {
			users = append(users, client)
		}
	}

	sort.Strings(users)

	return users
}

func broadcastToAll(response WsJSONResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket error", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

// renderPage renders the given template to the given ResponseWriter using the given data.
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
