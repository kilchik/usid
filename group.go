package main

import (
	"net/http"
)

type group struct {
	name    string
	users   map[*user]bool
	join    chan *user
	leave   chan *user
	forward chan message
}

func newGroup(name string) *group {
	return &group{
		name:    name,
		users:   make(map[*user]bool),
		join:    make(chan *user),
		leave:   make(chan *user),
		forward: make(chan message),
	}
}

func (g *group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logI.Println("in group", r.Cookies())
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logE.Printf("upgrade ws connection for group %q: %v", g.name, err)
	}

	cookie, err := r.Cookie("auth")
	if err != nil {
		logE.Printf("extract cookie: %v", err)
		return
	}

	userName := cookie.Value
	u := newUser(userName, conn, g)
	g.join <- u
	defer func() {
		g.leave <- u
	}()
	go u.write()
	u.read()
}

func (g *group) run() {
	for {
		select {
		case u := <-g.join:
			g.users[u] = true
			logI.Printf("user %s joined group %q", u.name, g.name)
		case u := <-g.leave:
			g.users[u] = true
			u.conn.Close()
			logI.Printf("user %s left group %q", u.name, g.name)
		case msg := <-g.forward:
			for u, online := range g.users {
				if online {
					u.send <- msg
				}
			}
		}
	}
}
