package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"

	"github.com/go-chi/chi/v5"
)

type ScheduleHandler struct{
	ServiceSchedule *service.ScheduleService
}

func (as *ScheduleHandler)AddSlot(w http.ResponseWriter, r *http.Request){
	var schedule domain.ScheduleSlot

	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := as.ServiceSchedule.AddSlot(schedule)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (gst *ScheduleHandler)GetScheduleToday(w http.ResponseWriter, r *http.Request){
	//schedule/today
	//GET /schedule/today?user_id=1&day_of_week=2&parity=even
	getUserID := r.URL.Query().Get("user_id")
	getDayOfWeek := r.URL.Query().Get("day_of_week")
	getParity := r.URL.Query().Get("parity")

	if getUserID == "" || getDayOfWeek == "" || getParity == ""{
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	intUserID, err := strconv.Atoi(getUserID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	intDayOfWeek, err := strconv.Atoi(getDayOfWeek)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := gst.ServiceSchedule.GetScheduleToday(intUserID, intDayOfWeek, domain.WeekParity(getParity))
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result == nil{
		http.Error(w, "404 NOT FOUND", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (uts *ScheduleHandler) UpdateTaskSchedule(w http.ResponseWriter, r *http.Request){
	//schedule/{id}
	var schedule domain.ScheduleSlot
	task_id := chi.URLParam(r, "id")
	if task_id == ""{
		http.Error(w, "err", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	intuserID , err := strconv.Atoi(task_id)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	schedule.ID = intuserID

	err = uts.ServiceSchedule.UpdateTaskSchedule(schedule)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}