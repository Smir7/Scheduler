package main

import (
	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/handler"
	"github.com/smir7/scheduler/repository"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

func main() {
	dataBase := repository.CreateDatabase()
	defer dataBase.Close()

	repo := repository.NewDatabase(dataBase)

	port := constans.Port
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) != 0 {
		port = envPort
	}
	port = ":" + port
	fileServer := http.FileServer(http.Dir(constans.WebDir))

	http.Handle("/", fileServer)
	http.HandleFunc("/api/nextdate", handler.NextDate)
	http.HandleFunc("GET /api/task", handler.TaskGet(repo))
	http.HandleFunc("POST /api/task", handler.TaskPost(repo))
	http.HandleFunc("PUT /api/task", handler.TaskPut(repo))
	http.HandleFunc("DELETE /api/task", handler.TaskDelete(repo))
	http.HandleFunc("/api/tasks", handler.TasksGet(repo))
	http.HandleFunc("/api/task/done", handler.TaskDone(repo))

	log.Println("Application is running", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
