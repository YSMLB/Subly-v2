package handler

import (
	"encoding/json"
	"net/http"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
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