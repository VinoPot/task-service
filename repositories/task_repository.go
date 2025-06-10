package repositories

import (
	"fmt"
	"task-service/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(db *gorm.DB) {
	DB = db
	DB.AutoMigrate(&models.Task{})
}

func GetAllTasks() []models.Task {
	if DB == nil {
		fmt.Println("❌ DB is nil!")
	} else {
		fmt.Println("✅ DB is initialized")
	}

	var tasks []models.Task
	DB.Find(&tasks)
	return tasks
}

func GetTaskByID(id string) models.Task {
	var task models.Task
	DB.First(&task, id)
	return task
}

func CreateTask(task *models.Task) {
	DB.Create(task)
}

func UpdateTask(id string, newTask *models.Task) {
	var task models.Task
	DB.First(&task, id)
	task.Title = newTask.Title
	task.Description = newTask.Description
	task.Done = newTask.Done
	DB.Save(&task)
}

func DeleteTask(id string) {
	var task models.Task
	DB.Delete(&task, id)
}
