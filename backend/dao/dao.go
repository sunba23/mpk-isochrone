package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sunba23/mpkIsoEngine/models"
)

func GetTravelData(stopId int) (map[int]models.TravelData, error) {
	var config Config = *GetConfig()
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	err = db.Ping()
	CheckError(err)

	statement := `
  SELECT to_stop_id, travel_time, route_stops
  FROM precomputed_travel_times
  WHERE from_stop_id = $1
  `

	rows, err := db.Query(statement, stopId)
	CheckError(err)
	defer rows.Close()

	travelDataMap := make(map[int]models.TravelData)

	for rows.Next() {
		var fromStopId, toStopId, travelTime int
		var routeStops []int

		err := rows.Scan(&toStopId, &travelTime, &routeStops)
		CheckError(err)

		travelDataMap[fromStopId] = models.TravelData{
			Id:         toStopId,
			TravelTime: travelTime,
			Path:       routeStops,
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating through rows: ", err)
		return nil, err
	}

	return travelDataMap, nil
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
