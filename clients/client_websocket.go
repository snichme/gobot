package clients

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/snichme/gobot/types"
)

type WebsocketClient struct {
	robot types.Robot
}

type WsConnection struct {
	writer http.ResponseWriter
	robot  types.Robot
	query  string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(robot types.Robot) http.HandlerFunc {
	isAuthorized := func(username, password string) bool {
		users := robot.Brain().Get("users").([]map[string]string)
		for _, user := range users {
			if user["username"] == username && user["password"] == password {
				return true
			}
		}
		return false
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if ok == false || !isAuthorized(username, password) {
			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		for {
			messageType, r, err := conn.NextReader()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			w, err := conn.NextWriter(messageType)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := io.Copy(w, r); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := w.Close(); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
}

func (client WebsocketClient) Start() {
	uri := "0.0.0.0:" + client.robot.Setting("websocket_port")
	fmt.Fprintf(client.robot, "WebsocketClient: Listening on ws://%s/ws\n", uri)
	http.HandleFunc("/ws", wsHandler(client.robot))
	http.ListenAndServe(uri, nil)
}

func NewWebsocketClient(robot types.Robot) *WebsocketClient {

	return &WebsocketClient{
		robot: robot,
	}
}
