package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
	"testing"
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
	return nil, nil
}

func (m *MockChangeRepo)UpdateStatusChanges(slotID int, userID int, newStatus domain.ChangeStatus) (*domain.ProposedChange, error){
	return nil, nil
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