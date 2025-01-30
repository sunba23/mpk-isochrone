package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sunba23/mpkIsoEngine/dao"
)

type travelDataResponse struct {
	code                int
	stopIdTravelDataMap map[int][]byte
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
		code:                http.StatusOK,
		stopIdTravelDataMap: stopTravelData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
