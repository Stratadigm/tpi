package tpi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stratadigm/tpi_data"
	"google.golang.org/appengine"
	"net/http"
	_ "net/http/httptest"
	"strconv"
	"testing"
	_ "time"
)

func TestCRUDUser(t *testing.T) {

	var err error
	g1 := &tpi_data.User{Name: "Roger", Email: "rafa@fed.com", Password: "allurbase", Confirmed: true, Rep: 0}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(g1)
	if err != nil {
		t.Errorf("Encode json : %v", err)
	}
	//Create
	//req, err := http.NewRequest("POST", "https://thalipriceindex.appspot.com/create/user", &buf)
	req, err := http.NewRequest("POST", "http://192.168.0.9:8080/create/user", &buf)
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

}

func TestCRUDVenue(t *testing.T) {

	var err error
	g1 := &tpi_data.Venue{Name: "Udupi Uphara", Location: appengine.GeoPoint{Lat: 13.9, Lng: 75.4}}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(g1)
	if err != nil {
		t.Errorf("Encode json : %v", err)
	}
	//Create
	req, err := http.NewRequest("POST", "https://thalipriceindex.appspot.com/create/venue", &buf)
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

func TestCRUDThali(t *testing.T) {

}

func TestCRUDData(t *testing.T) {

}

func TestLogin(t *testing.T) {

	var err error
	g1 := &tpi_data.User{Name: "Rafa", Email: "rafa@fed.com", Password: "allurbase"}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(g1)
	if err != nil {
		t.Errorf("Encode json : %v", err)
	}
	//Login
	req, err := http.NewRequest("POST", "http://192.168.0.9:8080/token_auth", &buf)
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response: %v", resp.Status)
	}
	cl := resp.Header.Get("Content-Length")
	icl, err := strconv.Atoi(cl)
	if err != nil {
		t.Errorf("Content Length err: %v", err)
	}
	tok := make([]byte, icl)
	dec := bufio.NewReader(resp.Body)
	_, err = dec.Read(tok)
	if err != nil {
		t.Errorf("Json resp err: %v", err)
	} //else {
	//	t.Errorf("Json token: %v", string(tok))
	//}
	resp.Body.Close()

	//Test token
	req, err = http.NewRequest("GET", "http://192.168.0.9:8080/hello", nil)
	if err != nil {
		t.Errorf("Request : %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", string(tok))) // Authorized
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", "abcd")) // Unauthorized
	//req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Client do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response: %v", resp.Status)
	}
	cl = resp.Header.Get("Content-Length")
	icl, err = strconv.Atoi(cl)
	if err != nil {
		t.Errorf("Content Length err: %v", err)
	}
	hw := make([]byte, icl)
	dec = bufio.NewReader(resp.Body)
	_, err = dec.Read(hw)
	if err != nil {
		t.Errorf("Json resp err: %v", err)
	}

}

func TestRefresh(t *testing.T) {

	var err error
	g1 := &tpi_data.User{Name: "Rafa", Email: "rafa@fed.com", Password: "allurbase"}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(g1)
	if err != nil {
		t.Errorf("Encode json : %v", err)
	}
	//Create
	//req, err := http.NewRequest("POST", "https://thalipriceindex.appspot.com/create/user", &buf)
	req, err := http.NewRequest("POST", "http://192.168.0.9:8080/refresh_token_auth", &buf)
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response: %v", resp.Status)
	}
	cl := resp.Header.Get("Content-Length")
	icl, err := strconv.Atoi(cl)
	if err != nil {
		t.Errorf("Content Length err: %v", err)
	}
	tok := make([]byte, icl)
	dec := bufio.NewReader(resp.Body)
	_, err = dec.Read(tok)
	if err != nil {
		t.Errorf("Json resp err: %v", err)
	} else {
		t.Errorf("Json token: %v", string(tok))
	}

}
