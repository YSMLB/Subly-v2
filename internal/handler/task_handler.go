package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"

	"github.com/go-chi/chi/v5"
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

func (gtid *TaskHandler) GetTaskID(w http.ResponseWriter, r *http.Request){
	getURLParametrs := chi.URLParam(r, "id")
	task, err := strconv.Atoi(getURLParametrs)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if task <= 0{
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	taskID, err := gtid.ServiceTask.GetTaskID(task)

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if taskID == nil{
		http.Error(w, "404 NOT FOUND", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskID)
}