package tpi

import (
	"github.com/stratadigm/tpi_settings"
	"google.golang.org/appengine"
	"net/http"
)

var (
//validEmail = regexp.MustCompile("^.*@.*\\.(com|org|in|mail|io)$")
)

func init() {

	if appengine.IsDevAppServer() {
		//fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		tpi_settings.SetEnvironment("preproduction")
	}
	tpi_settings.LoadSettingsByEnv(tpi_settings.GetEnvironment())

	r := NewRouter()
	http.Handle("/", r)

}
