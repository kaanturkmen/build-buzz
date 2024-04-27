package handler

import (
	"encoding/json"
	"github.com/kaanturkmen/build-buzz/internal/request"
	"github.com/kaanturkmen/build-buzz/internal/service"
	"log"
	"net/http"
)

type NotifyHandler struct {
	notifyService *service.NotifyService
}

func NewNotifyHandler(notifyService *service.NotifyService) *NotifyHandler {
	return &NotifyHandler{notifyService: notifyService}
}

func (h *NotifyHandler) Notify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse := struct {
			Error string `json:"error"`
		}{
			Error: "only POST method is allowed",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		jsonErr := json.NewEncoder(w).Encode(errorResponse)
		if jsonErr != nil {
			log.Printf("error encoding JSON: %v", jsonErr)
		}
		return
	}

	var notifyReq request.NotifyRequest
	err := json.NewDecoder(r.Body).Decode(&notifyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.notifyService.Notify(notifyReq)

	if err != nil {
		errorResponse := struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		jsonErr := json.NewEncoder(w).Encode(errorResponse)
		if jsonErr != nil {
			log.Printf("error encoding JSON: %v", jsonErr)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding JSON: %v", err)
	}
}
