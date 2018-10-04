package main

import (
	"github.com/gorilla/websocket"
	"time"
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
		logI.Printf("read message: %q", string(msg))
		logI.Printf("forward message: %v", message{u.name, string(msg), time.Now().Local()})
		u.group.forward <- message{u.name, string(msg), time.Now().Local()}
	}
}

func (u *user) write() {
	defer u.conn.Close()
	for msg := range u.send {
		logI.Printf("write %q to user %s", msg.Text, u.name)
		if err := u.conn.WriteJSON(&msg); err != nil {
			logE.Printf("write to user %s: %v", u.name, err)
			break
		}
	}
}
