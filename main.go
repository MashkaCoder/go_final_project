package main

import (
	"log"
	"net/http"
	"os"

	dbfuncs "github.com/MashkaCoder/go_final_project/database"
	"github.com/MashkaCoder/go_final_project/handler"
)

func main() {
	db, err := dbfuncs.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetConnection().Close()

	h := &handler.Handler{DB: db}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}
	log.Printf("Сервер слушает на порту %s", port)

	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)
	http.HandleFunc("/api/task", h.TaskHandler)
	http.HandleFunc("/api/tasks", h.TasksHandler)
	http.HandleFunc("/api/task/done", h.TaskDoneHandler)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}

}
