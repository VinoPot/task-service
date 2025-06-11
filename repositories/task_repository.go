package repositories

import (
	"task-service/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(db *gorm.DB) error {
	DB = db
	DB.AutoMigrate(&models.Task{})
	return nil
}

func GetAllTasks() []models.Task {
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

func UpdateTask(id string, newTask *models.Task) error {
	var task models.Task
	result := DB.First(&task, id)

	if result.Error != nil {
		return result.Error
	}

	task.Title = newTask.Title
	task.Description = newTask.Description
	task.Done = newTask.Done

	DB.Save(&task)
	return nil
}

func DeleteTask(id string) error {
	var task models.Task
	result := DB.Delete(&task, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
