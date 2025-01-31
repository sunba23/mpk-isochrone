package dao

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
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
		var toStopId, travelTime int
		var routeStops []int64

		err := rows.Scan(&toStopId, &travelTime, pq.Array(&routeStops))
		CheckError(err)

		travelDataMap[toStopId] = models.TravelData{
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

func GetStopsDetails() (map[string]models.StopData, error) {
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
  SELECT stop_id, stop_code, stop_name, ST_AsEWKB(stop_loc::geometry)
  FROM stops;
  `

	rows, err := db.Query(statement)
	CheckError(err)
	defer rows.Close()

	stopsDetailsMap := make(map[string]models.StopData)

	for rows.Next() {
		var stopId string
		var stopCode string
		var stopName string
		var stopLoc []byte

		err := rows.Scan(&stopId, &stopCode, &stopName, &stopLoc)
		CheckError(err)

		stopsDetailsMap[stopId] = models.StopData{
			Id:       stopId,
			Code:     stopCode,
			Name:     stopName,
			Location: stopLoc,
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating through rows: ", err)
		return nil, err
	}

	return stopsDetailsMap, nil
}
