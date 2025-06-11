package handlers

import (
	"net/http"
	"task-service/models"
	"task-service/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Глобальная переменная для подключения к БД
var DB *gorm.DB

func InitDB(db *gorm.DB) {
	repositories.Init(db)
	DB = db
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Description get a list of tasks with optional filtering by done status
// @ID get-tasks
// @Produce json
// @Param done query boolean false "Filter by completion status"
// @Success 200 {array} models.Task
// @Router /tasks [get]
func GetAllTasks(c *gin.Context) {
	done := c.Query("done")

	var tasks []models.Task

	if done == "true" {
		DB.Where("done = ?", true).Find(&tasks)
	} else if done == "false" {
		DB.Where("done = ?", false).Find(&tasks)
	} else {
		DB.Find(&tasks)
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Description get task by ID
// @ID get-task-by-id
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task := repositories.GetTaskByID(id)

	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask godoc
// @Summary Create a new task
// @Description create a new task with title and description
// @ID create-task
// @Accept json
// @Produce json
// @Param task body models.Task true "Task object"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repositories.CreateTask(&task)
	c.JSON(http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary Update an existing task
// @Description update task by ID
// @ID update-task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.Task true "Updated task object"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в данных задачи"})
		return
	}

	err := repositories.UpdateTask(id, &updatedTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description delete task by ID
// @ID delete-task
// @Produce json
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := repositories.DeleteTask(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.Status(http.StatusNoContent)
}
