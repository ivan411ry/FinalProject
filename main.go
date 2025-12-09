package main

import (
	"FinalProject/pkg/api"
	"FinalProject/pkg/db"

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

	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
