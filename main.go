package main

import (
	"task-service/config"
	"task-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()
	handlers.InitDB(db)

	r := gin.Default()

	r.GET("/tasks", handlers.GetAllTasks)
	r.GET("/tasks/:id", handlers.GetTaskByID)
	r.POST("/tasks", handlers.CreateTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.Run(":8080")
}
