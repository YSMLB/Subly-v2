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

type MockChangeRepo struct{}

func (m *MockChangeRepo)SaveChanges(changes domain.ProposedChange) (*domain.ProposedChange, error) {
	return &domain.ProposedChange{
		ID: 1,
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionCancel,
		Status: domain.ChangeStatusPending,
		RawText: "MDK 05.02",
		NewStartTime: "10:10",
		NewEndTime: "11:40",
		NewRoom: "306",
	}, nil
}

func (m *MockChangeRepo)GetPendingByUserID(userID int) ([]domain.ProposedChange, error){
	return nil, nil
}

func (m *MockChangeRepo)GetChangeByID(id int) (*domain.ProposedChange, error){
	return &domain.ProposedChange{
		ID: 1,
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "test",
		NewStartTime: "test",
		NewEndTime: "test",
		NewRoom: "test",
	}, nil
}

func (m *MockChangeRepo)UpdateStatusChanges(slotID int, userID int, newStatus domain.ChangeStatus) (*domain.ProposedChange, error){
	return &domain.ProposedChange{
		ID: 1,
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "test",
		NewStartTime: "test",
		NewEndTime: "test",
		NewRoom: "test",
	}, nil
}

func TestChangehandler_Save(t *testing.T){

	type TestCase struct{
		name string
		requestBody string
		expectedCode int
	}

	tests:= []TestCase{
		{
			name: "success: Save Changes",
			requestBody: `{"id": 1, "userID": 1, "slotID": 1,
			"changeAction": "cancel", "changeStatus": "pending",
			"rawText": "MDK 05.02", "newStartTime": "10:10",
			"newEndTime": "11:40", "newRoom": "306"
			}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "fail: broken JSON",
			requestBody: `{"ID: asd, "UserID": 1, "SlotID": 1,
			"Action": "pending", "Status": "52",
			"RawText": "MDK 05.02", "NewStartTime": "10:10",
			"NewEndTime": "11:40", "NewRoom": "306",
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "fail: business logic",
			requestBody: `{"ID: "asd", "UserID": 1, "SlotID": 1,
			"Action": "action", "Status": "52",
			"RawText": "MDK 05.02", "NewStartTime": 10,
			"NewEndTime": "11:40", "NewRoom": "306",
			}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	mockRepo := &MockChangeRepo{}
	changeService := &service.ChangeService{Repo: mockRepo}
	changeHandler := &ChangeHandler{ChangeService: changeService}

	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			req:= httptest.NewRequest(http.MethodPost, "/change", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			changeHandler.SaveChanges(rr, req)

			if rr.Code != tc.expectedCode{
				t.Errorf("Ожидался статус-код %d, получили %d", tc.expectedCode, rr.Code)
			}
		})
	}
}

func TestChangeHandler_GetChangeByID(t *testing.T){

	type TestCase struct{
		name string
		requestBody string
		url string
		expectedCode int
	}

	tests :=[]TestCase{
		{
			name: "success",
			url: "/changes/1",
			requestBody: `{"id": 1, "userID": 1, "slotID": 1,
			"changeAction": "cancel", "changeStatus": "pending",
			"rawText": "MDK 05.02", "newStartTime": "10:10",
			"newEndTime": "11:40", "newRoom": "306"
			}`,
			expectedCode: http.StatusOK,
		},
		{
			name: "fail: not natural",
			requestBody: `{"id": 1, "userID": 1, "slotID": 1,
			"changeAction": "cancel", "changeStatus": "pending",
			"rawText": "MDK 05.02", "newStartTime": "10:10",
			"newEndTime": "11:40", "newRoom": "306"
			}`,
			url: "/changes/asd",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "fail",
			requestBody: `{"id": 1, "userID": 1, "slotID": 1,
			"changeAction": "cancel", "changeStatus": "pending",
			"rawText": "MDK 05.02", "newStartTime": "10:10",
			"newEndTime": "11:40", "newRoom": "306"
			}`,
			url: "/changes/",
			expectedCode: http.StatusNotFound,
		},
	}

	mockRepo := &MockChangeRepo{}
	changeService := &service.ChangeService{Repo: mockRepo}
	changeHandler := &ChangeHandler{ChangeService: changeService}

	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rr := httptest.NewRecorder()

			test := chi.NewRouter()
			test.Get("/changes/{id}", changeHandler.GetChangeByID)
			test.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode{
				t.Errorf("ожидался статус-код %d, получили %d", tc.expectedCode, rr.Code)
			}
		})
	}
}


func TestChangeHandler_UpdateStatusChange(t *testing.T){

	type TestCase struct{
		name string
		url string
		requestBody string
		expectedCode int
	}

	tests := []TestCase{
		{
			name: "success",
			url: "/changes/1/status",
			requestBody: `{"id": 1, "userID": 1, "slotID": 1,
			"changeAction": "cancel", "changeStatus": "approved",
			"rawText": "MDK 05.02", "newStartTime": "10:10",
			"newEndTime": "11:40", "newRoom": "306"
			}`,
			expectedCode: http.StatusOK,
		},
		{
			name: "fail: id url",
			url: "/changes/abc/status",
			requestBody: `{"id": 2, "userID": 2, "slotID": 2,
			"changeAction": "cancel", "changeStatus": "pending",
			"rawText": "MDK 05.02", "newStartTime": "10:10",
			"newEndTime": "11:40", "newRoom": "306"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "fail: broken JSON",
			url: "/changes/3/status",
			requestBody: `{"id: 1, "userID": 1, "slotID": 1,
			"changeAction": "cancel", "changeStatus": "pending",
			"rawText": "MDK 05.02", "newStartTime": "10:10",
			"newEndTime": "11:40", "newRoom": "306"
			}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	mockRepo := &MockChangeRepo{}
	changeService := &service.ChangeService{Repo: mockRepo}
	changeHandler := &ChangeHandler{ChangeService: changeService}

	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, tc.url, bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			test := chi.NewRouter()
			test.Patch("/changes/{id}/status", changeHandler.UpdateStatusChanges)
			test.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode{
				t.Errorf("ожидался статус-код %d, пришел %d", tc.expectedCode, rr.Code)
			}
		})
	}
}

func TestChangeHandler_GetPendingByID(t *testing.T){

	type TestCase struct{
		name string
		requestBody string //change/pending?user_id=1
		expectedCode int
	}

	tests := []TestCase{
		{
			name: "success",
			requestBody: "/change/pending?user_id=1",
			expectedCode: http.StatusOK,
		},
		{
			name: "fail: broken JSON",
			requestBody: "/change/pending?user_id=asd",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "fail: user id <= 0",
			requestBody: "/change/pending?user_id=0",
			expectedCode: http.StatusBadRequest,
		},
	}

	//Важная херня!!
	mockRepo := &MockChangeRepo{}
	changeService := &service.ChangeService{Repo: mockRepo}
	changeHandler := &ChangeHandler{ChangeService: changeService}
	//Важная херня!!

	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.t){
			req := httptest.NewRequest(http.MethodGet, tc.requestBody, nil)
			rr := httptest.NewRecorder()

			test := chi.NewRouter()
			test.Get("/change/pending", changeHandler.GetPendingByUserID)
			test.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode{
				t.Errorf("ожидался статус-код %d, пришел %d", tc.expectedCode, rr.Code)
			}
		})
	}
}