package handlers

import (
	"net/http"
	"task-service/models"
	"task-service/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbInstance *gorm.DB) {
	repositories.Init(dbInstance)
}

func GetAllTasks(c *gin.Context) {
	tasks := repositories.GetAllTasks()
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task := repositories.GetTaskByID(id)
	done := c.Query("done")

	var tasks []models.Task

	if done == "true" {
		DB.Where("done = ?", true).Find(&tasks)
	} else if done == "false" {
		DB.Where("done = ?", false).Find(&tasks)
	} else {
		DB.Find(&tasks)
	}
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repositories.CreateTask(&task)
	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repositories.UpdateTask(id, &updatedTask)
	c.JSON(http.StatusOK, updatedTask)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	repositories.DeleteTask(id)
	c.Status(http.StatusNoContent)
}
