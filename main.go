package main

import (
	"FinalProject/pkg/api"
	"FinalProject/pkg/db"
	"log"

	"net/http"
	"os"
)

func main() {
	serv := os.Getenv("TODO_DBFILE")
	if serv == "" {
		serv = "scheduler.db"
	}
	err := db.Init(serv)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
