package models

type TravelData struct {
	Id         int
	TravelTime int
	Path       []int64
}

type StopData struct {
  Id int
  Code int
  Name string
  Location []byte
}
