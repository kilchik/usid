package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type user struct {
	conn  *websocket.Conn
	name  string
	send  chan message
	group *group
}

func newUser(name string, conn *websocket.Conn, g *group) *user {
	return &user{
		name:  name,
		conn:  conn,
		send:  make(chan message),
		group: g,
	}
}

func (u *user) read() {
	for {
		_, msg, err := u.conn.ReadMessage()
		if err != nil {
			logE.Printf("read input of user %s: %v", u.name, err)
			break
		}
		u.group.forward <- message{u, string(msg)}
	}
}

func (u *user) write() {
	defer u.conn.Close()
	for msg := range u.send {
		logI.Printf("write %q to user %s", msg.text, u.name)
		text := fmt.Sprintf("%s: %s", msg.from.name, msg.text)
		if err := u.conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			logE.Printf("write to user %s: %v", u.name, err)
			break
		}
	}
}
