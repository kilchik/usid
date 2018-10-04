package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/signature"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	th.templ.Execute(w, &data)
}

type message struct {
	From string    `json:"from"`
	Text string    `json:"text"`
	At   time.Time `json:"at"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	confPath := flag.String("config", "usid.conf", "path to configuration file")
	flag.Parse()

	conf, err := readConfig(*confPath)
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	logI = log.New(os.Stdout, "[I] ", log.Ldate|log.Ltime)
	logE = log.New(os.Stderr, "[E] ", log.Ldate|log.Ltime)

	logI.Println("started :)")

	g := newGroup("default")
	go g.run()

	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New(conf.GoogleClientId, conf.GoogleClientSecret,
			fmt.Sprintf("http://%s/auth/callback/google", conf.Listen)),
	)

	http.HandleFunc("/auth/", authHandler)
	http.Handle("/", &withAuth{&templateHandler{tname: "usid.html", g: g}})
	http.Handle("/login", &templateHandler{tname: "login.html", g: g})
	http.Handle("/group", g)

	logI.Printf("listening at %s", conf.Listen)
	if err := http.ListenAndServe(conf.Listen, nil); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
