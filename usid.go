package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	logE *log.Logger
	logI *log.Logger
)

type templateHandler struct {
	tname string
	once  sync.Once
	templ *template.Template
	g     *group
}

func (th templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.once.Do(func() {
		th.templ = template.Must(template.ParseFiles(filepath.Join("templates", th.tname)))
	})
	th.templ.Execute(w, nil)
}

type message struct {
	from *user
	text string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	logI = log.New(os.Stdout, "[I] ", log.Ldate|log.Ltime)
	logE = log.New(os.Stderr, "[E] ", log.Ldate|log.Ltime)

	g := newGroup("default")
	go g.run()
	http.Handle("/usid", &templateHandler{tname: "usid.html", g: g})
	http.Handle("/usid/group", g)
	if err := http.ListenAndServe("127.0.0.1:9632", nil); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
