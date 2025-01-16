package main

import (
	"log"
	"net/http"
	"os"

	
	dbfuncs "todo/database"
	"todo/handler"
	
)

func main(){
    err := dbfuncs.CreateTable()
    if err != nil{
        log.Fatal(err)
    }
    
	http.Handle("/", http.FileServer(http.Dir("web")))
    http.HandleFunc("/api/nextdate",  handler.NextDateHandler)
    http.HandleFunc("/api/task",      handler.TaskHandler)
    http.HandleFunc("/api/tasks",     handler.TasksHandler)
    http.HandleFunc("/api/task/done", handler.TaskDoneHandler)
    port := os.Getenv("TODO_PORT")
    if port == ""{
        port = ":7540" 
    }
    log.Printf("Сервер слушает на порту %s", port)

    defer dbfuncs.GetDB().Close()  
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal(err)
    }

}