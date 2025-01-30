package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sunba23/mpkIsoEngine/dao"
)

type travelDataResponse struct {
	Code                int            `json:"code"`
	StopIdTravelDataMap map[int][]byte `json:"stop_id_travel_data_map"`
}

type stopsDetailsResponse struct {
	Code           int               `json:"code"`
	StopDetailsMap map[string][]byte `json:"stop_details_map"`
}

func travelData(w http.ResponseWriter, r *http.Request) {
	stopIdStr := r.URL.Query().Get("stop_id")
	if stopIdStr == "" {
		http.Error(w, "Missing stop_id query parameter", http.StatusBadRequest)
		return
	}

	stopId, err := strconv.Atoi(stopIdStr)
	if err != nil {
		http.Error(w, "Invalid stop_id", http.StatusBadRequest)
		return
	}

	stopTravelData, err := dao.GetTravelData(stopId)

	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}

	response := travelDataResponse{
		Code:                http.StatusOK,
		StopIdTravelDataMap: stopTravelData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func stopsDetails(w http.ResponseWriter, _ *http.Request) {
	stopsDetails, err := dao.GetStopsDetails()

	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}

	response := stopsDetailsResponse{
		Code:           http.StatusOK,
		StopDetailsMap: stopsDetails,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
