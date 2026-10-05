package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
)

type MockTaskRepo1 struct{}

func (m *MockTaskRepo1) CreateTask(domain.Task) error{
	return nil
}

func (m *MockTaskRepo1) GetTaskID(int) (*domain.Task, error){
	return nil, nil
}

func (m *MockTaskRepo1) GetList(int) ([]domain.Task, error){
	return nil, nil
}

func (m *MockTaskRepo1) UpdateTask(*domain.Task) ( error){
	return nil
}

func (m *MockTaskRepo1) DeleteTask(int) error{
	return nil
}

func TestTaskHandler_Create(t *testing.T){
	//описываем из каких колонок состоит наша таблица
	type testCase struct{
		name string// имя кейса
		requestBody string//json входной
		expectedCode int //какой статус ждем
	}

	tests := []testCase{
		{
			name: "success: task created",
			requestBody: `{"title": "mockCreate Test", "id": 1, "userID": 1}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "fail: mistake from JSON file",
			requestBody: `{"title": "mockBadRequest Test", id: 2, "userID": 2}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "fail: mistake from businesses logic",
			requestBody: `{"title": "", "id": 0, "userID": 0}`,
			expectedCode: http.StatusBadRequest,
		},

	}

	//ВАЖНАЯ ХЕРНЯ!!!
	mockRepo := &MockTaskRepo1{}
	taskService := &service.TaskService{Repo: mockRepo}
	taskHandler := &TaskHandler{ServiceTask: taskService}
	//ВАЖНАЯ ХРЕНЬ!!!

	//регаем сам движок
	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			taskHandler.CreateTask(rr, req)

			if rr.Code != tc.expectedCode{
				t.Errorf("ожидался статус код %d, а получили %d", tc.expectedCode, rr.Code)
			}
		})
	}

}

func TestTaskHandler_GetTaskID(t *testing.T){
	
	//Определейние структуры (шаг1)
	type testCase struct{
		name string
		requestBody string
		expectedCode int
	}

	//Срез с набором проводимых сценариев (шаг2)
	tests := []testCase{
		{//успешыенй сценгарий
			name: "success: get task id",
			requestBody: "/tasks/1",
			expectedCode: http.StatusOK,
		},
		{//нет цифр, вместо них буквы
			name: "fail: field have not int",
			requestBody: "/tasks/sashaLox",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "fail: DB have not this user",
			requestBody: "/tasks/999",
			expectedCode: http.StatusNotFound,
		},
	}

	//ВАЖНАЯ ХРЕНЬТ!!!!!!!!
	mockRepo := &MockTaskRepo1{}
	taskService := &service.TaskService{Repo: mockRepo}
	taskHandler := &TaskHandler{ServiceTask: taskService}
	//РЯЛ ВАЖДНАЯ!!!

	//регаем движок
	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.requestBody, nil)
			rr := httptest.NewRecorder()

			test := chi.NewRouter()
			test.Get("/tasks/{id}", taskHandler.GetTaskID)
			test.ServeHTTP(rr, req)//вызов

			if rr.Code != tc.expectedCode{
				t.Errorf("ожидался статус-код %d, а получен %d", tc.expectedCode, rr.Code)
			}
		})
	}
}