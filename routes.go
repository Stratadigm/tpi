package tpi

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Create",
		"POST",
		"/create/user",
		Create,
	},
	Route{
		"Retrieve",
		"GET",
		"/retrieve",
		Retrieve,
	},
	Route{
		"Update",
		"POST",
		"/update",
		Update,
	},
	Route{
		"Delete",
		"POST",
		"/delete",
		Delete,
	},
	Route{
		"Logs",
		"GET",
		"/logs",
		Logs,
	},
	Route{
		"Users",
		"GET",
		"/users",
		Users,
	},
	Route{
		"Counters",
		"GET",
		"/counters",
		Counters,
	},
}
