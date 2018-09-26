package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/usid", mainHandler)
	if err := http.ListenAndServe("127.0.0.1:9632", nil); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

}
