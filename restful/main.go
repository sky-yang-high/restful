package main

import (
	taskserver "restful/taskserver"

	gin "github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	server := taskserver.NewTaskServer()

	router.POST("/task/create", server.CreateTaskHandler)
	router.GET("/task/all", server.GetAllTasksHandler)
	router.GET("/task/get/:id", server.GetTaskHandler)
	//router.GET("/task/tags/:tag", server.GetTaskByTagHandler)
	//router.PUT("/task/:id", server.UpdateTaskHandler)
	//router.DELETE("/task/:id", server.DeleteTaskHandler)
	router.Run(":8080")
}
