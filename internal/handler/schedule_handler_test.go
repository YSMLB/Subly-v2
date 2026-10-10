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

type MockScheduleRepo struct{}

func (m *MockScheduleRepo)AddSlot(schedule domain.ScheduleSlot) error{
	return nil
}

func (m *MockScheduleRepo)GetScheduleToday(userID int, dayOfWeek int) (*domain.ScheduleSlot, error){
	return &domain.ScheduleSlot{
		Parity: domain.WeekParityEven,
		UserID: 1,
		DayOfWeek: 2,
		Subject: "География",
		Teacher: "Иванова",
		Room: "302",
		StartTime: "10:10",
		EndTime: "11:40",
		
	}, nil
}

func (m *MockScheduleRepo)UpdateTaskSchedule(schedule domain.ScheduleSlot) error{
	return nil
}

func TestSchedule_AddSlot(t *testing.T){

	type TestCase struct{
		name string
		requestBody string
		expectedCode int
	}

	tests := []TestCase{
		{
			name: "fail: json",
			requestBody: `{"ID": 1, "UserID": 1, "Subject": "subject", "Teacher": "М.И.",
			"Room": "301", "DayOfWeek": 1, "StartTime": "10:10",
			"EndTime": "11:40",
			"Type": "lecture", "Parity": "even, 
			"IsCancelled": false}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "false: business rules",
			requestBody: `{"ID": 0, "UserID": -1, "Subject": "subject", "Teacher": "М.И.",
			"Room": "301", "DayOfWeek": 1, "StartTime": "10:10",
			"EndTime": "11:40",
			"Type": "lecture", "Parity": "even", 
			"IsCancelled": false}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "success: Add Slot",
			requestBody:`{"ID": 1, "UserID": 1, "Subject": "subject", "Teacher": "М.И.",
			"Room": "301", "DayOfWeek": 1, "StartTime": "10:10",
			"EndTime": "11:40",
			"Type": "lecture", "Parity": "even", 
			"IsCancelled": false}`,
			expectedCode: http.StatusCreated,
		},
	}

	mockRepo := &MockScheduleRepo{}
	scheduleService := &service.ScheduleService{Repo: mockRepo}
	scheduleHandler := &ScheduleHandler{ServiceSchedule: scheduleService}


	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			scheduleHandler.AddSlot(rr, req)

			if rr.Code != tc.expectedCode{
				t.Errorf("ожидался статус-код %d, получили %d", tc.expectedCode, rr.Code)
			}
		})
	}
}

func TestGetScheduleToday(t *testing.T){

	type TestCase struct{
		name string
		requestBody string
		expectedCode int
	}

	tests := []TestCase{
		{
			name: "success: get schedule today",
			requestBody: "/schedule/today?user_id=1&day_of_week=2&parity=even",
			expectedCode: http.StatusOK,
		},
		{
			name: "fail: broken data",
			requestBody: "/schedule/today?user_id=0&day_of_week=9&parity=...",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "fail: words",
			requestBody: "/schedule/today?user_id=asd&day_of_week=2&parity=even",
			expectedCode: http.StatusBadRequest,
		},
	}

	mockRepo := &MockScheduleRepo{}
	scheduleService := &service.ScheduleService{Repo: mockRepo}
	scheduleHandler := &ScheduleHandler{ServiceSchedule: scheduleService}

	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.requestBody, nil)
			rr := httptest.NewRecorder()

			test := chi.NewRouter()
			test.Get("/schedule/today", scheduleHandler.GetScheduleToday)
			
			test.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode{
				t.Errorf("ожидался стасус-код %d, получен %d", tc.expectedCode, rr.Code)
			}
		})
	}
}