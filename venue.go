package tpi

import (
	"google.golang.org/appengine/datastore"
	"time"
)

type Venue struct {
	Id        int64   `json:"id" schema:"-"`
	Name      string  `json:"name" schema:"name"`
	Latitude  float64 `json:"latitude" schema:"latitude"`
	Longitude float64 `json:"longitude" schema:"longitude"`
	Thalis    []Thali `json:"thalis"`
}

type VenueDatabase interface {
	ListVenues() ([]*Venue, error)

	AddVenue(guesty *Venue) (int64, error) //create

	GetVenue(id int64) (*Venue, error) //retrieve by id

	GetVenuewEmail(email string) (*Venue, error) //retrieve by email

	GetVenueKey(email string) (*Venue, *datastore.Key, error)

	UpdateVenue(guesty *Venue) error //update

	DeleteVenue(id int64) error //delete

	Close() error
}

func NewVenue(id int64) *Venue {

	return &Venue{Id: id}

}
