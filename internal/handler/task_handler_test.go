package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
	"testing"

//	"github.com/go-chi/chi/v5"
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