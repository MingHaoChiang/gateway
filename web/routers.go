package web

import (
	"fmt"
	"net/http"
	"strings"

	"com.hyweb/gateway/common"
	"github.com/gorilla/mux"
)

//Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes ...
type Routes []Route

//NewRouter ...
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := common.Logger(route.HandlerFunc, route.Name, route.Pattern)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

//Index ...
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html>gateway service</html>")
}

var routes = Routes{
	Route{
		"Index",
		strings.ToUpper("Get"),
		"/gateway/index.html",
		Index,
	},
	Route{
		"GetJwtToken",
		strings.ToUpper("Post"),
		"/jwt/token",
		GetJwtToken,
	},
	Route{
		"VerifyJWTToken",
		strings.ToUpper("Get"),
		"/jwt/token",
		VerifyJWTToken,
	},
	Route{
		"RenewJWTToken",
		strings.ToUpper("Put"),
		"/jwt/token",
		RenewJWTToken,
	},
}
