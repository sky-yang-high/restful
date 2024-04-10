package main

import (
	"restful/middleware"
	taskserver "restful/taskserver"

	gin "github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	server := taskserver.NewTaskServer()

	//since the method function didn't call the c.Next, the middleware should be registered before that.
	//more deeply, if a middleware is registered after the method function, it will not be executed.
	//that is because this middleware is not in the chain of the method function
	//(see the code like router.GET, and you will konw why)
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	router.POST("/task/create", server.CreateTaskHandler)
	router.GET("/task/all", server.GetAllTasksHandler)
	router.GET("/task/get/:id", server.GetTaskHandler)
	//router.GET("/task/tags/:tag", server.GetTaskByTagHandler)
	//router.PUT("/task/:id", server.UpdateTaskHandler)
	//router.DELETE("/task/:id", server.DeleteTaskHandler)
	router.GET("/panic", func(ctx *gin.Context) {
		panic("This is a panic")
	})

	//router.RunTLS()
	router.Run(":8080")
}
