package tpi

import (
	"google.golang.org/appengine/datastore"
	"image"
	"time"
)

type Data struct {
	TThali     Thali     `json:"thali"`
	TVenue     Venue     `json:"ven"`
	SubmitTime time.Time `json:"submitTime"`
	TUser      User      `json:"contributor"`
	Verfied    bool      `json:"verified"`
	Accepted   bool      `json:"accepted"`
}

type DataDatabase interface {
	ListDatas() ([]*Data, error)

	AddData(guesty *Data) (int64, error) //create

	GetData(id int64) (*Data, error) //retrieve by id

	GetDatawEmail(email string) (*Data, error) //retrieve by email

	GetDataKey(email string) (*Data, *datastore.Key, error)

	UpdateData(guesty *Data) error //update

	DeleteData(id int64) error //delete

	Close() error
}

func NewData(id int64) *Data {

	return &Data{Id: id}

}
