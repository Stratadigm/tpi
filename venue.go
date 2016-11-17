package tpi

import (
	"google.golang.org/appengine"
)

type Venue struct {
	Id   int64  `json:"id" schema:"-"`
	Name string `json:"name" schema:"name"`
	//Latitude  float64 `json:"latitude" schema:"latitude"`
	//Longitude float64 `json:"longitude" schema:"longitude"`
	Location appengine.GeoPoint `json:"location" schema:"location"`
	//Thalis   []Thali            `json:"thalis"`
	//Thalis []int64 `json:"thalis"`
}

func NewVenue(id int64) *Venue {

	return &Venue{Id: id}

}
