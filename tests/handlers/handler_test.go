package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"task-service/handlers"
	"task-service/models"
	"task-service/repositories"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func setupRouter(t *testing.T) *gin.Engine {
	if DB == nil {
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
		handlers.InitDB(DB)
	}

	DB.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")

	r := gin.Default()
	r.GET("/tasks", handlers.GetAllTasks)
	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks/:id", handlers.GetTaskByID)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	return r
}

func TestGetAllTasks(t *testing.T) {
	router := setupRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "[]")
}

func TestCreateTask_ValidationError(t *testing.T) {
	router := setupRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestGetTaskByID(t *testing.T) {
	router := setupRouter(t)

	taskJSON := `{"title": "Test", "description": "Description", "done": false}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Task
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	req2, _ := http.NewRequest("GET", "/tasks/"+strconv.Itoa(int(created.ID)), nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Test")
}

func TestUpdateTask(t *testing.T) {
	router := setupRouter(t)

	taskJSON := `{"title": "To Update", "description": "Old", "done": false}`
	req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var created models.Task
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	updateJSON := `{"title": "Updated", "description": "New", "done": true}`
	req2, _ := http.NewRequest("PUT", "/tasks/"+strconv.Itoa(int(created.ID)), strings.NewReader(updateJSON))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Updated")
}

func TestDeleteTask(t *testing.T) {
	router := setupRouter(t)

	taskJSON := `{"title": "To Delete", "description": "Delete Me", "done": false}`
	req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var created models.Task
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	req2, _ := http.NewRequest("DELETE", "/tasks/"+strconv.Itoa(int(created.ID)), nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNoContent, w2.Code)
}
