package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

const sqlGetRecentMessages = `
SELECT message_id, name, message, time, color
FROM (
	SELECT *
	FROM messages
	ORDER BY message_id DESC
	LIMIT 20) AS sub
ORDER BY message_id
`

var colorIndex = 0
var colors = []string{
	"#000",
	"#800",
	"#070",
	"#008",
	"#444",
	"#870",
	"#808",
	"#078",
	"#878",
}

func getColor() string {
	colorIndex += 1
	if colorIndex >= len(colors) {
		colorIndex = 0
	}

	return colors[colorIndex]
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
		addr:   req.RemoteAddr,
		color:  getColor(),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()

	go func() {
		rows, err := db.Query(sqlGetRecentMessages)
		if err != nil {
			return
		}

		for rows.Next() {
			var msg Message
			if err := rows.Scan(&msg.Id, &msg.Name, &msg.Message, &msg.Time, &msg.Color); err != nil {
				return
			}
			msg.Type = TextMessage

			res, err := json.Marshal(msg)
			if err != nil {
				return
			}

			client.send <- res
		}
	}()

	client.read()
}
