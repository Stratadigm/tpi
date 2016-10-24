package tpi

import (
	"bytes"
	"encoding/json"
	"net/http"
	_ "net/http/httptest"
	"testing"
	"time"
)

func TestCRUDUser(t *testing.T) {

	var err error
	g1 := &User{Name: "Roger", Email: "roger@fed.com", Confirmed: true, Points: nil, Rep: 0, JDte: time.Now()}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(g1)
	if err != nil {
		t.Errorf("Encode json : %v", err)
	}
	//Create
	req, err := http.NewRequest("POST", "https://thalipriceindex.appspot.com/create/user", &buf)
	if err != nil {
		t.Errorf("Request : %v", err)
	}
	//req.Header.Set("X-Custom-Header", "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Client do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Response: %v", resp.Status)
	}

	//Retrieve

	//g2 := &User{}
	//dec := json.NewDecoder(resp.Body)
	//err = json.Decode(g2)

	//Update

	//Delete

}

func TestCRUDVenue(t *testing.T) {

}

func TestCRUDThali(t *testing.T) {

}

func TestCRUDData(t *testing.T) {

}
