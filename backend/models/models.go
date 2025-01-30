package models

type TravelData struct {
	Id         int     `json:"id"`
	TravelTime int     `json:"travel_time"`
	Path       []int64 `json:"path"`
}

type StopData struct {
	Id       int
	Code     int
	Name     string
	Location []byte
}
