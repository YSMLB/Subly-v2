package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
)

type MockTaskRepo1 struct{}

func (m *MockTaskRepo1)CreateTask(domain.Task) error{
	return nil
}

func (m *MockTaskRepo1)GetTaskID(id int) (*domain.Task, error){
	return &domain.Task{
		ID: 1,
		UserID: 1,
		Title: "ананьев чмо",
	}, nil
}

func (m *MockTaskRepo1) GetList(userID int) ([]domain.Task, error){
	return nil, nil
}

func (m *MockTaskRepo1) UpdateTask(task *domain.Task) error{
	return nil
}

func (m *MockTaskRepo1) DeleteTask(taskID int) error{
	return nil
}

//тестовые сценарии для createtask
//позитивный сценарий
func TestTaskHandler_CreateTask_Success1(t *testing.T){
	//Подготовка (arange)
	mockRepo := &MockTaskRepo1{}//подключаем мок репу
	taskService := &service.TaskService{Repo: mockRepo}
	taskHandler := &TaskHandler{ServiceTask: taskService}

	//создаем мок json запрос
	requestBody := `{"userID": 1, "title": "афанасьев", "description": "лох"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(requestBody)) //создаем типо фейк запрос
	req.Header.Set("Content-Type", "application/json")//стандарт для json

	rr := httptest.NewRecorder()//starting test

	//calling handler (Act)
	taskHandler.CreateTask(rr, req)

	//check result (assert)
	if rr.Code != http.StatusCreated{
		t.Errorf("ожидался статус-код %d, получен %d", http.StatusCreated, rr.Code)
	}

	var createdTask domain.Task
	err := json.NewDecoder(rr.Body).Decode(&createdTask)
	if err != nil{
		t.Fatalf("Не удалось распарсить тело ответа: %v", err)
	}

	if createdTask.Title != "афанасьев"{
		t.Errorf("ожидался заголовок 'афанасьев', получен '%s'", createdTask.Title)
	}
}

func TestTaskHandler_CreateTask_InvalidJSON1(t *testing.T) {
	// Arrange
	mockRepo := &MockTaskRepo1{}
	taskService := &service.TaskService{Repo: mockRepo}
	taskHandler := &TaskHandler{ServiceTask: taskService}

	invalidBody := `{"title": "незакрытый json`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(invalidBody))
	rr := httptest.NewRecorder()

	// Act
	taskHandler.CreateTask(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("ожидался статус-код %d при битом JSON, получен %d", http.StatusBadRequest, rr.Code)
	}
}

// 3. Негативный сценарий: нарушение доменных правил (пустой заголовок)
func TestTaskHandler_CreateTask_ValidationError1(t *testing.T) {
	// Arrange
	mockRepo := &MockTaskRepo1{}
	taskService := &service.TaskService{Repo: mockRepo}
	taskHandler := &TaskHandler{ServiceTask: taskService}

	emptyTitleBody := `{"userID": 1, "title": ""}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(emptyTitleBody))
	rr := httptest.NewRecorder()

	// Act
	taskHandler.CreateTask(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("ожидался статус-код %d при нарушении валидации заголовка, получен %d", http.StatusBadRequest, rr.Code)
	}
}


func TestTaskHandler_GetTaskID_Success2(t *testing.T) {
	mockRepo := &MockTaskRepo1{}
	taskService := &service.TaskService{Repo: mockRepo}
	taskHandler := &TaskHandler{ServiceTask: taskService}

	requestBody := `{"ID": 1, "userID": 1, "title": "Тест2", "description": "nah"}`
	req := httptest.NewRequest(http.MethodGet, "/tasks/2", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()
	
	test := chi.NewRouter()
	test.Get("/tasks/{id}", taskHandler.GetTaskID)
	test.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK{
		t.Errorf("ожидался статус-код %d, получен %d", http.StatusOK, rr.Code)
	}

	var createdTask domain.Task
	err := json.NewDecoder(rr.Body).Decode(&createdTask)
	if err != nil{
		t.Fatalf("Не удалось распарсить тело ответа: %v", err)
	}
	if createdTask.ID != 1{
		t.Error("err")
	}

}
