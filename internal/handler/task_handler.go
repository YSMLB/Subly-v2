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


func (gl *TaskHandler) GetList(w http.ResponseWriter, r *http.Request){
	userid := r.URL.Query().Get("user_id")
	if userid == ""{
		http.Error(w, "notfound", http.StatusBadRequest)//400
		return
	}
	
	task, err := strconv.Atoi(userid)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)//400
		return
	}

	finally, err := gl.ServiceTask.GetList(task)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)//500
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finally)
}

func (ut *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request){
	getURLParametres := chi.URLParam(r, "id")// url - /tasks/{id}/status
	id, err:= strconv.Atoi(getURLParametres)//конвертим и получаем user id

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id <= 0{
		http.Error(w, "dont correct id", http.StatusBadRequest)
		return
	}

	type ChangeStruct struct{
		Userid int
		Status domain.TaskStatus
	}

	var change ChangeStruct//получаем новый статус
	err = json.NewDecoder(r.Body).Decode(&change)

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	finally, err := ut.ServiceTask.UpdateStatus(id, change.Userid, change.Status)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finally)
}

func  (dt *TaskHandler)DeleteTask(w http.ResponseWriter, r *http.Request){
	// задача лежит /tasks/{id}
	//идентификатор владельца соответственно в ?user_id=...

	taskID := chi.URLParam(r, "id")
	userID := r.URL.Query().Get("user_id")

	if userID == ""{
		http.Error(w, "not user_id", http.StatusBadRequest)
		return
	}

	if taskID == ""{
		http.Error(w, "not taskID", http.StatusBadRequest)
		return
	}

	TaskId, err := strconv.Atoi(taskID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} 

	UserId, err := strconv.Atoi(userID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dt.ServiceTask.DeleteTask(TaskId, UserId)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}