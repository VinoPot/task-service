// tests/repositories/repository_test.go
package repositories_test

import (
	"testing"

	"task-service/models"
	"task-service/repositories"

	"strconv"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func setupTestDB(t *testing.T) {
	dsn := "host=localhost user=postgres password=qaz123! dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Irkutsk"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("Ошибка миграции: %v", err)
	}

	repositories.Init(DB)

	// Очистка таблицы перед каждым тестом
	DB.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
}

func TestCreateTask(t *testing.T) {
	setupTestDB(t)

	task := &models.Task{
		Title:       "Test задача",
		Description: "Описание",
		Done:        false,
	}

	repositories.CreateTask(task)

	var savedTask models.Task
	DB.First(&savedTask, task.ID)

	assert.NotZero(t, savedTask.ID)
	assert.Equal(t, "Test задача", savedTask.Title)
	assert.False(t, savedTask.Done)
}

func TestGetTaskByID(t *testing.T) {
	setupTestDB(t)

	// Сначала создадим задачу
	task := &models.Task{
		Title:       "Test задача",
		Description: "Описание",
		Done:        false,
	}
	repositories.CreateTask(task)

	// Проверим GetTaskByID
	saved := repositories.GetTaskByID(strconv.Itoa(int(task.ID)))
	assert.Equal(t, "Test задача", saved.Title)
	assert.Equal(t, uint(task.ID), saved.ID)
}

func TestUpdateTask(t *testing.T) {
	setupTestDB(t)

	// Создаем задачу
	task := &models.Task{
		Title:       "Старое название",
		Description: "Описание",
		Done:        false,
	}
	repositories.CreateTask(task)

	// Обновляем её
	updated := &models.Task{
		ID:          task.ID,
		Title:       "Новое название",
		Description: "Обновлённое описание",
		Done:        true,
	}
	repositories.UpdateTask(strconv.Itoa(int(task.ID)), updated)

	var result models.Task
	DB.First(&result, task.ID)

	assert.Equal(t, "Новое название", result.Title)
	assert.Equal(t, "Обновлённое описание", result.Description)
	assert.True(t, result.Done)
}

func TestDeleteTask(t *testing.T) {
	setupTestDB(t)

	// Создаем задачу
	task := &models.Task{
		Title: "Задача для удаления",
		Done:  false,
	}
	repositories.CreateTask(task)

	// Удаляем её
	err := repositories.DeleteTask(strconv.Itoa(int(task.ID)))
	assert.NoError(t, err)

	// Проверяем, что задача действительно удалена
	var deletedTask models.Task
	result := DB.First(&deletedTask, task.ID)
	assert.Error(t, result.Error) // Должна быть ошибка (не найдено)
}

func TestGetAllTasks(t *testing.T) {
	setupTestDB(t)

	// Создаем несколько задач
	task1 := &models.Task{Title: "Задача 1", Description: "Первая задача"}
	task2 := &models.Task{Title: "Задача 2", Description: "Вторая задача"}

	repositories.CreateTask(task1)
	repositories.CreateTask(task2)

	tasks := repositories.GetAllTasks()
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Задача 1", tasks[0].Title)
	assert.Equal(t, "Задача 2", tasks[1].Title)
}
