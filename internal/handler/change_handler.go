package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"

	"github.com/go-chi/chi/v5"
)

type ChangeHandler struct{
	ChangeService *service.ChangeService
}

func (sc *ChangeHandler)SaveChanges(w http.ResponseWriter, r *http.Request){
	var change domain.ProposedChange

	err := json.NewDecoder(r.Body).Decode(&change)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := sc.ChangeService.SaveChanges(change)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (gpduid *ChangeHandler) GetPendingByUserID(w http.ResponseWriter, r *http.Request){
	// change/pending?user_id=1

	changeUserID := r.URL.Query().Get("user_id")

	intUserID , err:= strconv.Atoi(changeUserID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} 

	result, err := gpduid.ChangeService.GetPendingByUserID(intUserID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (gcbid *ChangeHandler)GetChangeByID(w http.ResponseWriter, r *http.Request){
	/// changes/{id}

	userid := chi.URLParam(r, "id")

	intUserID, err := strconv.Atoi(userid)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := gcbid.ChangeService.GetChangeByID(intUserID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}


func (usc *ChangeHandler)UpdateStatusChanges(w http.ResponseWriter, r *http.Request){
	// changes/{id}/status

	taskid := chi.URLParam(r, "id")
	intTaskID, err := strconv.Atoi(taskid)
	
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type status struct{
		UserID int
		Pending domain.ChangeStatus
	}

	var statusChange status
	err = json.NewDecoder(r.Body).Decode(&statusChange)

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := usc.ChangeService.UpdateStatusChanges(intTaskID, statusChange.UserID, statusChange.Pending)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}