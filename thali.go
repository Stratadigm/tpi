package tpi

import (
	"google.golang.org/appengine/datastore"
	"image"
	"time"
)

type Thali struct {
	Target  int // 1-4 target customer profile
	Limited bool
	Region  int     // 1-3 target cuisine
	Price   float64 //
	Photo   image
}

type ThaliDatabase interface {
	ListThalis() ([]*Thali, error)

	AddThali(guesty *Thali) (int64, error) //create

	GetThali(id int64) (*Thali, error) //retrieve by id

	GetThaliwEmail(email string) (*Thali, error) //retrieve by email

	GetThaliKey(email string) (*Thali, *datastore.Key, error)

	UpdateThali(guesty *Thali) error //update

	DeleteThali(id int64) error //delete

	Close() error
}

func NewThali(id int64) *Thali {

	return &Thali{Id: id}

}
