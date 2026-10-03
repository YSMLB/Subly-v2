package handler

import (
	"encoding/json"
	"net/http"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
)

type TaskHandler struct{
	ServiceTask *service.TaskService
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request){
	var task domain.Task//инициализируем запрос в переменной
	err := json.NewDecoder(r.Body).Decode(&task)//декодируем запрос

	if err != nil{//обрабатываем ошибку
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.ServiceTask.CreateTask(task)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)//обработка ошибки и возврат ее
		return
	}

	w.WriteHeader(http.StatusCreated)//отдаем код 201
	json.NewEncoder(w).Encode(res)//кодируем и отправляем тело JSON
}