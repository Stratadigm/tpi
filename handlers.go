package tpi

import (
	_ "appengine"
	_ "bytes"
	"encoding/base64"
	"encoding/json"
	_ "fmt"
	_ "github.com/gorilla/schema"
	_ "golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	_ "image/jpeg"
	"net/http"
	"time"
)

var (
	tmpl_logs = template.Must(template.ParseFiles("templates/logs"))
)

const recordsPerPage = 10

type Render struct { //for most purposes
	Average float64 `json:"average"`
}

// Index writes in JSON format the average value of a thali at the requester's location to the response writer
func Index(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	host := GetIp(r)
	loc, err := GetLoc(c, host)
	if err != nil {
		log.Errorf(c, "Index get location: %v", err)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(loc); err != nil {
		log.Errorf(c, "Index json encode: %v", err)
		return
	}
	return

}

//Create uses data in JSON post to create a User/Venue/Thali/Data
func Create(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	var g1 interface{}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	adsc := &DS{ctx: c}
	//Need to make sure counter is alive before creating/adding guests
	counter := adsc.GetCounter()
	if counter == nil {
		err := adsc.CreateCounter()
		if err != nil {
			log.Errorf(c, "Create create counter: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&DSErr{time.Now(), "Create create counter " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	}
	switch r.URL.Path {
	case "/create/user":
		g1 = &User{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create user: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&DSErr{time.Now(), "Create user " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/venue":
		g1 = &Venue{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create venue: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/thali":
		g1 = &Thali{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create thali: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&DSErr{time.Now(), "Create thali " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode DSErr: %v", err)
				return
			}
			return
		}
	case "/create/data":
		g1 = &Data{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create data: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&DSErr{time.Now(), "Create data " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode DSErr: %v", err)
				return
			}
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		if err := enc.Encode(&DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode DSErr: %v", err)
			return
		}
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(g1)
	if err != nil {
		log.Errorf(c, "Couldn't decode posted json: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := enc.Encode(&DSErr{time.Now(), "Create entity " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode DSErr: %v", err)
			return
		}
		return
	}
	if id, err := adsc.Add(g1); err != nil {
		log.Errorf(c, "Couldn't add entity: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := enc.Encode(&DSErr{time.Now(), "Create entity " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode DSErr: %v", err)
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		if err := enc.Encode(&DSErr{time.Now(), "Created entity " + string(id)}); err != nil {
			log.Errorf(c, "Created json encode DSErr: %v", err)
			return
		}
		return
	}
}

//Create uses data in posted form to create a User/Venue/Thali/Data
/*func Create(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	_ = r.ParseForm()
	var g1 interface{}
	enc := json.NewEncoder(w)
	adsc := &DS{ctx: c}
	//Need to make sure counter is alive before creating/adding guests
	counter := adsc.GetCounter()
	if counter == nil {
		err := adsc.CreateCounter()
		if err != nil {
			log.Errorf(c, "Create create counter: %v", err)
			if err := enc.Encode(&DSErr{time.Now(), "Create create counter " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	}
	switch r.URL.Path {
	case "/create/user":
		g1 = &User{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create user: %v", err)
			if err := enc.Encode(&DSErr{time.Now(), "Create user " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/venue":
		g1 = &Venue{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create venue: %v", err)
			if err := enc.Encode(&DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/thali":
		g1 = &Thali{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create thali: %v", err)
			if err := enc.Encode(&DSErr{time.Now(), "Create thali " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/data":
		g1 = &Data{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create data: %v", err)
			if err := enc.Encode(&DSErr{time.Now(), "Create data " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	default:
		if err := enc.Encode(&DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode: %v", err)
			return
		}
		return
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(g1, r.PostForm)
	if err != nil {
		log.Errorf(c, "Couldn't decode posted form: %v\n", err)
		return
	}
	if id, err := adsc.Add(g1); err != nil {
		log.Errorf(c, "Couldn't add entity: %v\n", err)
		if err := enc.Encode(&DSErr{time.Now(), "Create entity " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode: %v", err)
			return
		}
		return
	} else {
		if err := enc.Encode(&DSErr{time.Now(), "Created entity " + string(id)}); err != nil {
			log.Errorf(c, "Created json encode: %v", err)
			return
		}
		return
	}
}*/

//Retrieve gets the posted entity from datastore
func Retrieve(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	log.Errorf(c, "Retrieve")
	return

}

//Update updates the posted entity in datastore
func Update(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	log.Errorf(c, "Update")
	return

}

//Delete deletes the posted entity from datastore
func Delete(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	log.Errorf(c, "Delete")
	return

}

//Logs writes logs in html to the response writer
func Logs(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	var data struct {
		Records []*log.Record
		Offset  string
	}

	query := &log.Query{AppLogs: true}

	if offset := r.FormValue("offset"); offset != "" {
		query.Offset, _ = base64.URLEncoding.DecodeString(offset)
	}

	res := query.Run(ctx)

	for i := 0; i < recordsPerPage; i++ {
		rec, err := res.Next()
		if err == log.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "Reading log records: %v", err)
			break
		}

		data.Records = append(data.Records, rec)
		if i == recordsPerPage-1 {
			data.Offset = base64.URLEncoding.EncodeToString(rec.Offset)
		}
	}

	if err := tmpl_logs.Execute(w, data); err != nil {
		log.Errorf(ctx, "Rendering template: %v", err)
	}

}
