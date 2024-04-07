package main

import (
	"net/http"
	taskserver "restful/taskserver"

	mux "github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	server := taskserver.NewTaskServer()

	router.HandleFunc("/task/create", server.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/task/all", server.GetAllTasksHandler).Methods("GET")
	router.HandleFunc("/task/get/{id:[0-9]+}", server.GetTaskHandler).Methods("GET")
	//router.HandleFunc("/task/tags/{tag}", server.GetTaskByTagHandler).Methods("GET")
	//router.HandleFunc("PUT /task/{id}", server.UpdateTaskHandler).Methods("PUT")
	//router.HandleFunc("DELETE /task/{id}", server.DeleteTaskHandler).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
