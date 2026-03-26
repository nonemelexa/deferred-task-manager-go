package api

import (
	"log"
	"net/http"
)

func StartServer(handler *Handler, port string) {
	addr := ":" + port

	http.HandleFunc("/tasks", handler.CreateTask)

	log.Println("Server listening on", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
