package tpi

import (
	_ "image"
	_ "image/jpeg"
	"time"
)

type Thali struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Submitted time.Time `json:"submitted"`
	Target    string    `json:"target" schema:"target"` // 1-4 target customer profile
	Limited   bool      `json:"limited" schema:"limited"`
	Region    string    `json:"region" schema:"region"` // 1-3 target cuisine
	Price     float64   `json:"price" schema:"price"`
	//Photo     *image.RGBA `json:"image"`
	Photo    string `json:"image"`
	UserId   int64  `json:"userid"`
	VenueId  int64  `json:"venue" schema:"venue"`
	Verfied  bool   `json:"verified"`
	Accepted bool   `json:"accepted"`
}

func NewThali(id int64) *Thali {

	return &Thali{Id: id}

}
