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
		"Login",
		"POST",
		"/token_auth",
		Login,
	},
	/*Route{
		"Refresh",
		"GET",
		"/refresh_token_auth",
		RefreshToken,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		Logout,
	},*/
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Create",
		"POST",
		"/create/{what}",
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
		"List",
		"GET",
		"/list/{what}",
		List,
	},
	Route{
		"JSONList",
		"GET",
		"/jsonlist/{what}",
		JSONList,
	},
	Route{
		"Counters",
		"GET",
		"/counters",
		Counters,
	},
	Route{
		"PostForm",
		"POST",
		"/postform/{what}",
		PostForm,
	},
	Route{
		"GetForm",
		"GET",
		"/getform/{what}",
		GetForm,
	},
	Route{
		"GetFormId",
		"GET",
		"/getform/{what}/{id:[0-9]+}",
		GetForm,
	},
	Route{
		"GetUpload",
		"GET",
		"/upload/{what}",
		GetUpload,
	},
	Route{
		"PostUpload",
		"POST",
		"/upload/{what}",
		PostUpload,
	},
	Route{
		"GetImage",
		"GET",
		"/image/{what}",
		GetImage,
	},
}
