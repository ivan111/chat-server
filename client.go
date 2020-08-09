package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"strconv"
	"time"
	"unicode/utf8"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
	addr   string
	color  string
}

type MessageType int

const (
	TextMessage MessageType = iota + 1
	ErrorMessage
	InfoMessage
)

type Message struct {
	Type    MessageType `json:"type"`
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Time    time.Time   `json:"time"`
	Color   string      `json:"color"`
}

const sqlAddMessage = `
INSERT INTO messages(name, message, addr, color)
VALUES($1, $2, $3, $4)
RETURNING message_id
`

const sqlGetMessage = `
SELECT time
FROM messages
WHERE message_id = $1
`

func (c *client) read() {
	for {
		if _, bytes, err := c.socket.ReadMessage(); err == nil {
			var msg Message

			if err := json.Unmarshal(bytes, &msg); err != nil {
				c.sendErrorMessage("JSONの解読に失敗")
				continue
			}

			msg.Type = TextMessage
			msg.Color = c.color

			if msg.Name == "" {
				msg.Name = "Anonymous"
			}

			if utf8.RuneCountInString(msg.Name) > 16 {
				c.sendErrorMessage("名前が長すぎ。最大16文字まで")
				continue
			}

			if utf8.RuneCountInString(msg.Message) > 140 {
				c.sendErrorMessage("メッセージが長すぎ。最大140文字まで")
				continue
			}

			var idStr string
			err := db.QueryRow(sqlAddMessage, msg.Name, msg.Message, c.addr, c.color).Scan(&idStr)
			if err != nil {
				c.sendErrorMessage("データベースへの追加に失敗")
				continue
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.sendErrorMessage("IDの変換に失敗")
				continue
			}
			msg.Id = id

			var tm time.Time
			err = db.QueryRow(sqlGetMessage, idStr).Scan(&tm)
			if err != nil {
				c.sendErrorMessage("時刻の取得に失敗")
				continue
			}
			msg.Time = tm

			res, err := json.Marshal(msg)
			if err != nil {
				c.sendErrorMessage("JSONへの変換に失敗")
				continue
			}

			c.room.forward <- res

			m := normilizeMessage(msg.Message)
			if f, ok := responseFuncs[m]; ok {
				f(&msg, c.room)
			}
		} else {
			break
		}
	}

	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}

	c.socket.Close()
}

func (c *client) sendErrorMessage(msgText string) {
	var msg Message

	msg.Type = ErrorMessage
	msg.Message = msgText

	res, err := json.Marshal(msg)
	if err != nil {
		return
	}

	c.send <- res
}

func sendMessage(myRoom *room, name string, msgText string, color string) {
	var msg Message

	msg.Type = TextMessage
	msg.Color = color
	msg.Name = name
	msg.Message = msgText

	var idStr string
	err := db.QueryRow(sqlAddMessage, msg.Name, msg.Message, "localhost", color).Scan(&idStr)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}
	msg.Id = id

	var tm time.Time
	err = db.QueryRow(sqlGetMessage, idStr).Scan(&tm)
	if err != nil {
		return
	}
	msg.Time = tm

	res, err := json.Marshal(msg)
	if err != nil {
		return
	}

	myRoom.forward <- res
}
