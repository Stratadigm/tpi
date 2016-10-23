package tpi

import (
	"time"
)

type Data struct {
	Id         int64     `json:"id"`
	TThali     Thali     `json:"thali"`
	TVenue     Venue     `json:"ven"`
	SubmitTime time.Time `json:"submitTime"`
	UserId     int64     `json:"userid"`
	Verfied    bool      `json:"verified"`
	Accepted   bool      `json:"accepted"`
}

func NewData(id int64) *Data {

	return &Data{Id: id}

}
