package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tour-service/internal/service"

	"github.com/gorilla/mux"
)

type TourExecutionHandler struct {
	service *service.TourExecutionService
}

func NewTourExecutionHandler(service *service.TourExecutionService) *TourExecutionHandler {
	return &TourExecutionHandler{service: service}
}

func (h *TourExecutionHandler) StartTour(w http.ResponseWriter, r *http.Request) {
	touristID, _ := r.Context().Value("userID").(uint)
	
	vars := mux.Vars(r)
	tourID, _ := strconv.ParseUint(vars["tourId"], 10, 32)

	var req struct {
		StartLat float64 `json:"startLat"`
		StartLng float64 `json:"startLng"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	execution, err := h.service.StartTour(uint(tourID), touristID, req.StartLat, req.StartLng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(execution)
}

func (h *TourExecutionHandler) CheckPosition(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    executionID, _ := strconv.ParseUint(vars["executionId"], 10, 32)

    var req struct {
        CurrentLat float64 `json:"currentLat"`
        CurrentLng float64 `json:"currentLng"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    completed, err := h.service.CheckPosition(uint(executionID), req.CurrentLat, req.CurrentLng)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(completed)
}

func (h *TourExecutionHandler) CompleteTour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	executionID, _ := strconv.ParseUint(vars["executionId"], 10, 32)

	err := h.service.CompleteTour(uint(executionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) AbandonTour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	executionID, _ := strconv.ParseUint(vars["executionId"], 10, 32)

	err := h.service.AbandonTour(uint(executionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) GetActiveExecution(w http.ResponseWriter, r *http.Request) {
    touristID, _ := r.Context().Value("userID").(uint)
    
    vars := mux.Vars(r)
    tourID, _ := strconv.ParseUint(vars["tourId"], 10, 32)

    execution, err := h.service.GetActiveExecution(touristID, uint(tourID))
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(nil)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(execution) 
}

func (h *TourExecutionHandler) GetExecutionDetails(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    executionID, _ := strconv.ParseUint(vars["executionId"], 10, 32)

    execution, err := h.service.GetExecutionDetails(uint(executionID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(execution)
}

func (h *TourExecutionHandler) GetExecutionsByTour(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tourID, _ := strconv.ParseUint(vars["tourId"], 10, 32)

    executions, err := h.service.GetExecutionsByTour(uint(tourID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(executions)
}