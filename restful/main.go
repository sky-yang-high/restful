package main

import (
	"net/http"
	taskserver "restful/taskServer"
)

func main() {
	mux := http.NewServeMux()
	server := taskserver.NewTaskServer()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("GET /task/create", server.CreateTaskHandler)
	mux.HandleFunc("/task/all", server.GetAllTasksHandler)
	mux.HandleFunc("/task/get", server.GetTaskHandler)
	//mux.HandleFunc("/task/tags",server.GetTaskByTagHandler())
	//mux.HandleFunc("PUT /task/{id}", server.UpdateTaskHandler())
	//mux.HandleFunc("DELETE /task/{id}", server.DeleteTaskHandler())
	http.ListenAndServe(":8080", mux)
}
