package models

type TravelData struct {
	Id         int     `json:"id"`
	TravelTime int     `json:"travel_time"`
	Path       []int64 `json:"path"`
}

type StopData struct {
  Id       string `json:"id"`
  Code     string `json:"stop_code"`
  Name     string `json:"stop_name"`
  Location []byte `json:"stop_location"`
}
