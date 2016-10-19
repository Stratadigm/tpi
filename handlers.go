package tpi

import (
	_ "appengine"
	"bytes"
	"encoding/base64"
	"fmt"
	_ "golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"image/jpeg"
	"net/http"
	"time"
)

var (
	tmpl_cmn  = template.Must(template.ParseFiles("templates/base", "templates/head", "templates/body"))
	tmpl_logs = template.Must(template.ParseFiles("templates/logs"))
	tmpl_vnts = template.Must(template.ParseFiles("templates/events"))
	tmpl_err  = template.Must(template.ParseFiles("templates/base", "templates/head", "templates/err_body"))
	jsonFile  = "client_secret_250861196641-kk2ji01qp5ofa0hkml8uc8lda68ns4a6.apps.googleusercontent.com.json"
	tokenFile = "blackoutmap_token"
	filenames = []string{"f0.jpg", "f1.jpg", "f2.jpg"}
)

const recordsPerPage = 10

type Render struct { //for most purposes
	Message string   `json:"message"`
	Images  []string `json:"images"`
}

// Index writes in JSON format the average value of a thali at the requester's location to the response writer
func Index(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	images := make([]string, 0)
	for _, f := range filenames {
		buffer := new(bytes.Buffer)
		//b, err := ioutil.ReadFile(f) // for dev_appserver testing only
		img, err := ReadCloudImage(c, f) //ReadCloudImage (*image.Image, error)
		if err != nil {
			log.Errorf(c, "error reading from gcs %v \n", err)
			tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": f})
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
			tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": f})
			return
		}
		str := base64.StdEncoding.EncodeToString(buffer.Bytes())
		images = append(images, str) //buffer.Bytes())
	}
	data := Render{Message: "", Images: images}
	err := tmpl_cmn.ExecuteTemplate(w, "base", data)
	//w.Header().Set("Content-Type", "image/jpeg")
	//w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	//_, err := w.Write(buffer.Bytes())
	if err != nil {
		log.Errorf(c, "Couldn't execute common template: %v\n", err)
		tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": ""})
	}
	return

}

// Create uses data in JSON post to create a Data
func Create(w http.ResponseWriter, r *http.Request) {

	//var client *http.Client
	var url string
	c := appengine.NewContext(r)
	//first try the service account calendar
	if vnts, err := SAMoonPhases(c); err == nil && len(*vnts) != 0 {
		if err := tmpl_vnts.Execute(w, *vnts); err != nil {
			log.Errorf(c, "Rendering template: %v", err)
		}
		return
	}
	//then check for existing oauth2 token
	existing, err := ReadCloudToken(c, tokenFile)
	if existing != nil && err == nil && existing.AccessToken != "" && time.Now().Before(existing.Expiry) {
		client := cfg.Client(c, existing)
		events, err := MoonPhases(c, client)
		if err != nil {
			fmt.Fprint(w, "Sorry there was an error")
		}
		log.Debugf(c, "Number of events: %v", len(*events))
		if err := tmpl_vnts.Execute(w, *events); err != nil {
			log.Errorf(c, "Rendering template: %v", err)
		}
	} else { // else redirect user to google auth with authcode url
		scopes := []string{"https://www.googleapis.com/auth/calendar"}
		cfg, url = AuthCodeURL(jsonFile, scopes...)
		http.Redirect(w, r, url, http.StatusFound)
	}

	//after redirecting read from channel that paksha writes to and create client/token
	/*log.Debugf(c, "Getting token")
	  client, tok := OauthClientToken(c, cfg)
	  log.Debugf(c, "Got token: %v", tok)
	  if client != nil {
	          MoonPhases(client)

	  }
	  err = WriteCloudToken(c, tok, tokenFile)
	*/
	return

}

// Retrieve
func Retrieve(w http.ResponseWriter, r *http.Request) {

	/*c := appengine.NewContext(r)
	        existing, err := ReadCloudToken(c, tokenFile)
		if existing != nil && err == nil && existing.AccessToken != "" {
			client := cfg.Client(c, existing)
	                events, err := MoonPhases(c, client)
	                if err != nil {
	                        fmt.Fprint(w, "Sorry there was an error")
	                }
	                log.Debugf(c, "Number of events: %v", len(*events))
	                if err := tmpl_vnts.Execute(w, events); err != nil {
	                        log.Errorf(c, "Rendering template: %v", err)
	                }
		} else {
	                fmt.Fprint(w, "Sorry there was an error - please start over")
	        }
	        return*/

	c := appengine.NewContext(r)
	//x := daystogo()
	//data := Render{strconv.Itoa(x)+" days to go",}
	err := CreateImages(c)
	if err != nil {
		log.Errorf(c, "error creating image: %v\n", err)
		tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": ""})
	}

	images := make([]string, 0)
	for _, f := range filenames {
		buffer := new(bytes.Buffer)
		//b, err := ioutil.ReadFile(f) // for dev_appserver testing only
		img, err := ReadCloudImage(c, f) //ReadCloudImage (*image.Image, error)
		if err != nil {
			log.Errorf(c, "error reading from gcs %v \n", err)
			tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": f})
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
			tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": f})
			return
		}
		str := base64.StdEncoding.EncodeToString(buffer.Bytes())
		images = append(images, str) //buffer.Bytes())
	}
	data := Render{Message: "", Images: images}
	err = tmpl_cmn.ExecuteTemplate(w, "base", data)
	//w.Header().Set("Content-Type", "image/jpeg")
	//w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	//_, err := w.Write(buffer.Bytes())
	if err != nil {
		log.Errorf(c, "Couldn't execute common template: %v\n", err)
		tmpl_err.ExecuteTemplate(w, "base", map[string]interface{}{"Message": err, "Filename": ""})
	}
	return

}

//Update
func Update(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Welcome")
	return

}

//Deletes the posted
func Delete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.URL.Path == "/favicon.ico" {
		http.Error(w, "", 404)
		return
	}
	if r.FormValue("state") != randState {
		log.Debugf(c, "State doesn't match: req")
		http.Error(w, "", 500)
		return
	}
	var code string
	if code = r.FormValue("code"); code != "" {
		log.Debugf(c, "Got redirected to paksha")
		//w.(http.Flusher).Flush()
		//ch <- code
		client, tok := OauthClientToken(c, cfg, code)
		log.Debugf(c, "Got token: %v", tok)
		if client != nil {
			MoonPhases(c, client)

		}
		err := WriteCloudToken(c, tok, tokenFile)
		if err != nil {
			log.Debugf(c, "Writing token to gcs %v ", err)
		}
		log.Debugf(c, "Got code - authorized %v ", code)
	}
	log.Debugf(c, "Redirecting to moon")
	http.Redirect(w, r, "/moon", http.StatusFound)
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
