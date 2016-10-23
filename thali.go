package tpi

import (
	"image"
	_ "image/jpeg"
)

type Thali struct {
	Id      int64         `json:"id"`
	Target  int           `json:"target" schema:"target"` // 1-4 target customer profile
	Limited bool          `json:"limited" schema:"limited"`
	Region  int           `json:"region" schema:"region"` // 1-3 target cuisine
	Price   float64       `json:"price" schema:"price"`
	Photo   image.NRGBA64 `json:"image" schema:"image"`
}

func NewThali(id int64) *Thali {

	return &Thali{Id: id}

}
