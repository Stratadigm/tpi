package tpi

import (
	"github.com/gorilla/mux"
	"github.com/stratadigm/tpi_auth"
	"github.com/urfave/negroni"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile(`^/(create|jsonlist|list|users|counters|postform|getform|image|logs)?/?(.*)$`)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = makeHandler(route.HandlerFunc)
		//handler = Logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	router.Handle("/refresh_token_auth",
		negroni.New(
			negroni.HandlerFunc(tpi_auth.RequireTokenAuthentication),
			negroni.HandlerFunc(RefreshToken),
		)).Methods("GET")
	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(tpi_auth.RequireTokenAuthentication),
			negroni.HandlerFunc(Logout),
		)).Methods("GET")
	router.Handle("/hello",
		negroni.New(
			negroni.HandlerFunc(tpi_auth.RequireTokenAuthentication),
			negroni.HandlerFunc(Hello),
		)).Methods("GET")
	return router
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}
