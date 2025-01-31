package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sunba23/mpkIsoEngine/cache"
	"github.com/sunba23/mpkIsoEngine/dao"
	"github.com/sunba23/mpkIsoEngine/models"
)

type travelDataResponse struct {
	Code                int                       `json:"code"`
	StopIdTravelDataMap map[int]models.TravelData `json:"stop_id_travel_data_map"`
}

type stopsDetailsResponse struct {
	Code           int                        `json:"code"`
	StopDetailsMap map[string]models.StopData `json:"stop_details_map"`
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

	ctx := context.Background()
	rdb := cache.GetClient()

	cacheKey := fmt.Sprintf("travel_data:%s", stopIdStr)

	var stopTravelDataMap map[int]models.TravelData

  // try to get response from cache
	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &stopTravelDataMap); err != nil {
			log.Printf("Error unmarshaling cached data: %v", err)
		} else {
      log.Printf("Using value for key %v.", cacheKey)
			response := travelDataResponse{
				Code:                http.StatusOK,
				StopIdTravelDataMap: stopTravelDataMap,
			}
			sendResponse(w, response)
			return
		}
	} else if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

  // cache not found, so get travel travelData
	stopTravelDataMap, err = dao.GetTravelData(stopId)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

  // set cached value
  log.Printf("Setting value for key %v in redis at %v.", cacheKey, time.Now())
	if jsonData, err := json.Marshal(stopTravelDataMap); err == nil {
		err = rdb.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Error caching data: %v", err)
		}
	}

  log.Printf("Responding for key %v.", cacheKey)
	response := travelDataResponse{
		Code:                http.StatusOK,
		StopIdTravelDataMap: stopTravelDataMap,
	}
	sendResponse(w, response)
}

func stopsDetails(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	rdb := cache.GetClient()

	cacheKey := fmt.Sprintf("all_stops_details")

	var stopsDetailsMap map[string]models.StopData

	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &stopsDetailsMap); err != nil {
			log.Printf("Error unmarshaling cached data: %v", err)
		} else {
      log.Printf("Responding with cached value for key %v.", cacheKey)
			response := stopsDetailsResponse{
				Code:           200,
				StopDetailsMap: stopsDetailsMap,
			}
			sendResponse(w, response)
			return
		}
	} else if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	stopsDetailsMap, err = dao.GetStopsDetails()
	if err != nil {
    log.Printf("Database error: %s", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}

  log.Printf("Setting value for key %v in redis at %v.", cacheKey, time.Now())
	if jsonData, err := json.Marshal(stopsDetailsMap); err == nil {
		err = rdb.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Error caching data: %v", err)
		}
	}

  log.Printf("Responding for key %v.", cacheKey)
	response := stopsDetailsResponse{
		Code:           http.StatusOK,
		StopDetailsMap: stopsDetailsMap,
	}
  sendResponse(w, response)
}

func sendResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
