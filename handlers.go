package tpi

import (
	_ "appengine"
	"bytes"
	"encoding/base64"
	"encoding/json"
	_ "fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	//_ "golang.org/x/oauth2"
	"github.com/stratadigm/tpi_data"
	"github.com/stratadigm/tpi_services"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"html/template"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

var (
	tmpl_err        = template.Must(template.ParseFiles("templates/error"))
	tmpl_logs       = template.Must(template.ParseFiles("templates/logs"))
	tmpl_users      = template.Must(template.ParseFiles("templates/users"))
	tmpl_venues     = template.Must(template.ParseFiles("templates/venues"))
	tmpl_thalis     = template.Must(template.ParseFiles("templates/thalis"))
	tmpl_datas      = template.Must(template.ParseFiles("templates/datas"))
	tmpl_cntrs      = template.Must(template.ParseFiles("templates/counters"))
	tmpl_userform   = template.Must(template.ParseFiles("templates/cmn/base", "templates/cmn/body", "templates/userform"))
	tmpl_thaliform  = template.Must(template.ParseFiles("templates/cmn/base", "templates/cmn/body", "templates/thaliform"))
	tmpl_uploadform = template.Must(template.ParseFiles("templates/cmn/base", "templates/cmn/body", "templates/uploadform"))
	tmpl_image      = template.Must(template.ParseFiles("templates/cmn/base", "templates/cmn/body", "templates/image"))
)

const thanksMessage = `Thanks for input.`
const recordsPerPage = 10

//const perPage = 20

type Render struct { //for most purposes
	Average float64 `json:"average"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	requestUser := new(tpi_data.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := tpi_services.Login(c, requestUser)
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	c := appengine.NewContext(r)
	requestUser := new(tpi_data.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(tpi_services.RefreshToken(c, requestUser))

}

func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	c := appengine.NewContext(r)
	err := tpi_services.Logout(c, r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func Hello(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//c := appengine.NewContext(r)
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))

}

// Index writes in JSON format the average value of a thali at the requester's location to the response writer
func Index(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	host := tpi_data.GetIp(r)
	loc, err := tpi_data.GetLoc(c, host)
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

//Create uses data in JSON post to create a User/Venue/Thali/Data. Create first creates an empty entity & updates counter, then fills in fields using posted data and finally persists in datastore.
func Create(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	var g1 interface{}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	adsc := tpi_data.NewDSwc(c) //&DS{Ctx: c}
	//Need to make sure counter is alive before creating/adding entities
	counter := adsc.GetCounter()
	if counter == nil {
		err := adsc.CreateCounter()
		if err != nil {
			log.Errorf(c, "Create create counter: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create create counter " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	}
	switch r.URL.Path {
	case "/create/user":
		g1 = &tpi_data.User{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create user: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create user " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/venue":
		g1 = &tpi_data.Venue{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create venue: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/thali":
		g1 = &tpi_data.Thali{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create thali: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create thali " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode tpi_data.DSErr: %v", err)
				return
			}
			return
		}
	case "/create/data":
		g1 = &tpi_data.Data{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create data: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create data " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode tpi_data.DSErr: %v", err)
				return
			}
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode tpi_data.DSErr: %v", err)
			return
		}
		return
	}
	temp := reflect.ValueOf(g1).Elem().FieldByName("Id").Int()
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(g1)
	log.Errorf(c, "Creating: %v", g1)
	if err != nil {
		log.Errorf(c, "Couldn't decode posted json: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create entity decode " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode tpi_data.DSErr: %v", err)
			return
		}
		return
	}
	if err := adsc.Validate(g1); err != nil {
		log.Errorf(c, "Create json validate : %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create entity validation " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode tpi_data.DSErr: %v", err)
			return
		}
		return
	}
	//Need to specify Id when adding to datastore because json.Decode posted user data wipes out Id information
	reflect.ValueOf(g1).Elem().FieldByName("Id").SetInt(temp)
	reflect.ValueOf(g1).Elem().FieldByName("Submitted").Set(reflect.ValueOf(time.Now()))
	if id, err := adsc.Add(g1, temp); err != nil {
		log.Errorf(c, "Couldn't add entity: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create entity " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode tpi_data.DSErr: %v", err)
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Created entity " + string(id)}); err != nil {
			log.Errorf(c, "Created json encode tpi_data.DSErr: %v", err)
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
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create create counter " + err.Error()}); err != nil {
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
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create user " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/venue":
		g1 = &Venue{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create venue: %v", err)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/thali":
		g1 = &Thali{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create thali: %v", err)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create thali " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	case "/create/data":
		g1 = &Data{}
		if err = adsc.Create(g1); err != nil {
			log.Errorf(c, "Create data: %v", err)
			if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create data " + err.Error()}); err != nil {
				log.Errorf(c, "Create json encode: %v", err)
				return
			}
			return
		}
	default:
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create venue " + err.Error()}); err != nil {
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
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Create entity " + err.Error()}); err != nil {
			log.Errorf(c, "Create json encode: %v", err)
			return
		}
		return
	} else {
		if err := enc.Encode(&tpi_data.DSErr{time.Now(), "Created entity " + string(id)}); err != nil {
			log.Errorf(c, "Created json encode: %v", err)
			return
		}
		return
	}
}*/

//Retrieve gets list of entities of posted type from datastore
func Retrieve(w http.ResponseWriter, r *http.Request) {

	//var err error
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

//List writes list of Users/Venues/Thalis/Data in html to the response writer
func List(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	adsc := tpi_data.NewDSwc(c) //&DS{Ctx: c}
	var err error
	data := map[string]interface{}{
		"Next": "0",
		"Prev": "0",
	}

	offint := 0
	if offset := r.FormValue("offset"); offset != "" {
		offint, err = strconv.Atoi(offset)
		if err != nil {
			log.Errorf(c, "Reading user records offset: %v", err)
		}
	}

	//var g1 interface{}
	switch r.URL.Path {
	case "/list/users":
		g1 := make([]tpi_data.User, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "List users: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			tmpl_err.Execute(w, map[string]interface{}{"Message": err})
			return
		}
		data["List"] = g1
	case "/list/venues":
		g1 := make([]tpi_data.Venue, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "List venues: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			tmpl_err.Execute(w, map[string]interface{}{"Message": err})
			return
		}
		data["List"] = g1
	case "/list/thalis":
		g1 := make([]tpi_data.Thali, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "List thalis: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			tmpl_err.Execute(w, map[string]interface{}{"Message": err})
			return
		}
		data["List"] = g1
	case "/list/datas":
		g1 := make([]tpi_data.Data, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "List data: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			tmpl_err.Execute(w, map[string]interface{}{"Message": err})
			return
		}
		data["List"] = g1
	default:
		w.WriteHeader(http.StatusBadRequest)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad path"})
		return
	}

	data["Next"] = strconv.Itoa(offint + tpi_data.PerPage)
	if offint == 0 {
		data["Prev"] = strconv.Itoa(offint)
	} else {
		data["Prev"] = strconv.Itoa(offint - tpi_data.PerPage)
	}

	switch r.URL.Path {
	case "/list/users":
		if err := tmpl_users.Execute(w, data); err != nil {
			log.Errorf(c, "Rendering template: %v", err)
		}
	case "/list/venues":
		if err := tmpl_venues.Execute(w, data); err != nil {
			log.Errorf(c, "Rendering template: %v", err)
		}
	case "/list/thalis":
		if err := tmpl_thalis.Execute(w, data); err != nil {
			log.Errorf(c, "Rendering template: %v", err)
		}
	case "/list/datas":
		if err := tmpl_datas.Execute(w, data); err != nil {
			log.Errorf(c, "Rendering template: %v", err)
		}
	default:
		if err := tmpl_users.Execute(w, data); err != nil {
			log.Errorf(c, "Rendering template: %v", err)
		}
	}
}

//JSONList writes list of Users/Venues/Thalis/Data in html to the response writer
func JSONList(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	adsc := tpi_data.NewDSwc(c) //&DS{Ctx: c}
	var err error

	offint := 0
	if offset := r.FormValue("offset"); offset != "" {
		offint, err = strconv.Atoi(offset)
		if err != nil {
			log.Errorf(c, "Reading user records offset: %v", err)
		}
	}

	switch r.URL.Path {
	case "/jsonlist/users":
		g1 := make([]tpi_data.User, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "json list users: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err1 := enc.Encode(&tpi_data.DSErr{time.Now(), "Error " + err.Error()}); err1 != nil {
				log.Errorf(c, "json list users json encode err: %v", err1)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := enc.Encode(g1); err != nil {
			log.Errorf(c, "json list users encode: %v", err)
		}
		return
	case "/jsonlist/venues":
		g1 := make([]tpi_data.Venue, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "json list venues: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err1 := enc.Encode(&tpi_data.DSErr{time.Now(), "Error " + err.Error()}); err1 != nil {
				log.Errorf(c, "json list venues json encode err: %v", err1)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := enc.Encode(&g1); err != nil {
			log.Errorf(c, "json list venues encode: %v", err)
		}
		return
	case "/jsonlist/thalis":
		g1 := make([]tpi_data.Thali, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "json list thalis: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err1 := enc.Encode(&tpi_data.DSErr{time.Now(), "Error " + err.Error()}); err1 != nil {
				log.Errorf(c, "json list thalis json encode err: %v", err1)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := enc.Encode(&g1); err != nil {
			log.Errorf(c, "json list thalis encode: %v", err)
		}
		return
	case "/jsonlist/datas":
		g1 := make([]tpi_data.Data, 1)
		if err = adsc.List(&g1, offint); err != nil {
			log.Errorf(c, "json list data: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			if err1 := enc.Encode(&tpi_data.DSErr{time.Now(), "Error " + err.Error()}); err1 != nil {
				log.Errorf(c, "json list datas json encode err: %v", err1)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := enc.Encode(&g1); err != nil {
			log.Errorf(c, "json list datas encode: %v", err)
		}
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		if err1 := enc.Encode(&tpi_data.DSErr{time.Now(), "Error " + err.Error()}); err1 != nil {
			log.Errorf(c, "json list bad path err: %v", err1)
		}
		return
	}

}

//Users writes list of Users in html to the response writer
func Users(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	var err error
	var data struct {
		List []*tpi_data.User
		Next string
		Prev string
	}

	query := datastore.NewQuery("user").Order("Id")

	offint := 0
	if offset := r.FormValue("offset"); offset != "" {
		offint, err = strconv.Atoi(offset)
		if err != nil {
			log.Errorf(ctx, "Reading user records offset: %v", err)
		}
		query = query.Limit(tpi_data.PerPage + offint).Offset(offint)
	} else {
		query = query.Limit(tpi_data.PerPage).Offset(0)
	}

	users := make([]*tpi_data.User, 0)
	_, err = query.GetAll(ctx, &users)
	if err != nil {
		log.Errorf(ctx, "Datastore getall query: %v", err)
	}

	data.List = users
	data.Next = strconv.Itoa(offint + tpi_data.PerPage)
	if offint == 0 {
		data.Prev = strconv.Itoa(offint)
	} else {
		data.Prev = strconv.Itoa(offint - tpi_data.PerPage)
	}

	if err := tmpl_users.Execute(w, data); err != nil {
		log.Errorf(ctx, "Rendering template: %v", err)
	}

}

//Counters writes counter details in html to the response writer
func Counters(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	var err error

	query := datastore.NewQuery("counter")

	cntr := make([]*tpi_data.Counter, 0)
	_, err = query.GetAll(ctx, &cntr)
	if err != nil {
		log.Errorf(ctx, "Datastore getall query: %v", err)
	}

	if err := tmpl_cntrs.Execute(w, cntr[0]); err != nil {
		log.Errorf(ctx, "Rendering template: %v", err)
	}

}

//PostForm handles Post requests to create entities as specified in url path
func PostForm(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	_ = r.ParseForm()
	adsc := tpi_data.NewDSwc(c) //&DS{Ctx: c}
	//Need to make sure counter is alive before creating/adding guests
	counter := adsc.GetCounter()
	if counter == nil {
		err := adsc.CreateCounter()
		if err != nil {
			log.Errorf(c, "PostForm Create counter: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			tmpl_err.Execute(w, map[string]interface{}{"Message": "Couldn't create counter: " + err.Error()})
			return
		}
	}
	var g1 interface{}
	vars := mux.Vars(r)
	switch vars["what"] {
	case "user":
		g1 = &tpi_data.User{}
	case "venue":
		g1 = &tpi_data.Venue{}
	case "thali":
		g1 = &tpi_data.Thali{}
	case "data":
		g1 = &tpi_data.Data{}
	default:
	}
	if err = adsc.Create(g1); err != nil {
		log.Errorf(c, "PostForm Create : %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postform Create Error: " + err.Error()})
		return
	}
	temp := reflect.ValueOf(g1).Elem().FieldByName("Id").Int()
	decoder := schema.NewDecoder()
	err = decoder.Decode(g1, r.PostForm)
	if err != nil {
		log.Errorf(c, "Couldn't decode posted form: %v\n", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": err})
		return
	}
	if err := adsc.Validate(g1); err != nil {
		log.Errorf(c, "PostForm validate : %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postform Validate Error: " + err.Error()})
		return
	}
	//Need to specify Id when adding to datastore because schema.Decode posted user data wipes out Id information
	reflect.ValueOf(g1).Elem().FieldByName("Id").SetInt(temp)
	reflect.ValueOf(g1).Elem().FieldByName("Submitted").Set(reflect.ValueOf(time.Now()))
	if _, err = adsc.Add(g1, temp); err != nil {
		log.Errorf(c, "Postform add : %v\n", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postform Add Error: " + err.Error()})
		return
	}
	tmpl_err.Execute(w, map[string]interface{}{"Message": thanksMessage})
	return

}

//GetForm handles Get request to /getform/{what} and renders data input templates
func GetForm(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	switch vars["what"] {
	case "user":
		if err = tmpl_userform.ExecuteTemplate(w, "base", map[string]interface{}{"Message": thanksMessage}); err != nil {
			tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad get user form : " + err.Error()})
			return
		}
		return
	case "venue":
		if err = tmpl_thaliform.ExecuteTemplate(w, "base", map[string]interface{}{"Message": thanksMessage}); err != nil {
			tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad get venue form : " + err.Error()})
			return
		}
		return
	case "thali":
		var id string
		if _, ok := vars["id"]; !ok {
			id = "0"
		} else {
			id = vars["id"]
		}
		if err = tmpl_thaliform.ExecuteTemplate(w, "base", map[string]interface{}{"Id": id}); err != nil {
			tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad get thali form : " + err.Error()})
			return
		}
		return
	case "data":
		if err = tmpl_thaliform.ExecuteTemplate(w, "base", map[string]interface{}{"Message": thanksMessage}); err != nil {
			tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad get data form : " + err.Error()})
			return
		}
		return
	default:
		log.Errorf(c, "Bad getform url: %v", vars["what"])
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad getform url: " + vars["what"]})
		return

	}

}

//PostUpload handles Post requests to upload mulitpart files
func PostUpload(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["what"])
	if err != nil {
		log.Errorf(c, "Postupload strconv: %v", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postupload strconv: " + err.Error()})
	}
	adsc := tpi_data.NewDSwc(c) //&DS{Ctx: c}
	thali := &tpi_data.Thali{Id: int64(id)}
	if err = adsc.Get(thali); err != nil {
		log.Errorf(c, "Postupload get thali: %v", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postupload get thali: " + err.Error()})
		return
	}
	//_4MB := (1 << 17) * 4
	var file multipart.File
	//var header *multipart.FileHeader
	file, _, err = r.FormFile("image")
	if err != nil {
		log.Errorf(c, "Postupload formfile: %v", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postupload formfile: " + err.Error()})
		return
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf(c, "Postupload ReadAll: %v", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postupload ReadAll: " + err.Error()})
		return
	}
	rdr := bytes.NewReader(bs)
	img, _, err := image.Decode(rdr)
	if err != nil {
		log.Errorf(c, "Postupload Image decode: %v", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postupload Image decode: " + err.Error()})
		return
	}
	if err = tpi_data.WriteCloudImage(c, &img, vars["what"]); err != nil {
		log.Errorf(c, "Postupload Image write: %v", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Postupload Image write: " + err.Error()})
		return
	}
	thali.Photo = vars["what"]
	if _, err = adsc.Add(thali, thali.Id); err != nil {
		log.Errorf(c, "PostUpload Add : %v\n", err)
		tmpl_err.Execute(w, map[string]interface{}{"Message": "PostUpload Add : " + err.Error()})
		return
	}
	if _, ok := img.(*image.RGBA); ok {
		log.Errorf(c, "Postupload Image rgba: %v", ok)
	}
	tmpl_err.Execute(w, map[string]interface{}{"Message": "Postupload Success!!"})
	return

}

//GetUpload handles Get requests to file/image upload forms
func GetUpload(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["what"]
	if err = tmpl_uploadform.ExecuteTemplate(w, "base", map[string]interface{}{"Id": id}); err != nil {
		log.Errorf(c, "Bad getupload url: %v", vars["what"])
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad get upload form : " + err.Error()})
		return
	}
	return

}

//GetImage handles Get requests for images from GC bucket
func GetImage(w http.ResponseWriter, r *http.Request) {

	var err error
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["what"]

	buffer := new(bytes.Buffer)
	//b, err := ioutil.ReadFile(f) // for dev_appserver testing only
	img, err := tpi_data.ReadCloudImage(c, id) //ReadCloudImage (*image.Image, error)
	if err != nil {
		log.Errorf(c, "error reading from gcs %v \n", err)
		tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": id})
		return
	}
	//img, err := jpeg.Decode(bytes.NewReader(b)) //for dev_appserver testing only
	//if err != nil { //testing only
	//        log.Printf("error reading from gcs %v \n", err)
	//        tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message":err, "Filename":f})
	//        return
	//}//for dev_appserver testing only
	if err := jpeg.Encode(buffer, *img, nil); err != nil { //change *img to img for dev_appserver testing
		log.Errorf(c, "error reading image from gcs %v \n", err)
		tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": id})
		return
	}
	str := base64.StdEncoding.EncodeToString(buffer.Bytes())

	if err = tmpl_image.ExecuteTemplate(w, "base", map[string]interface{}{"Image": str}); err != nil {
		log.Errorf(c, "Bad image url: %v", vars["what"])
		tmpl_err.Execute(w, map[string]interface{}{"Message": "Bad image url : " + err.Error()})
		return
	}
	return

}
