package main

import (
	"log"
	"net/http"

	"github.com/andrewesteves/superchat/internal/handlers"
	"golang.org/x/net/websocket"
)

func main() {
	log.Println("app is running")

	talksHandler := handlers.NewTalks()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(talksHandler.Handle))

	fileServer := http.FileServer(http.Dir("./public/"))
	mux.Handle("/public/", http.StripPrefix("/public", fileServer))

	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
}
