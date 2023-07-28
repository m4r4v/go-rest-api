package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/m4r4v/go-rest-api/handlers"
	resources "github.com/m4r4v/go-rest-api/resources"
)

var data = &ServerData{
	apiVersion: "/v1",
	port:       "8080",
}

func ServerRouter() {

	// New Router Instance
	router := mux.NewRouter().StrictSlash(true)

	// Set Headers to accept only JSON requests
	// TODO
	// Show Content Type error message in JSON

	// Handle Error 404
	router.NotFoundHandler = http.HandlerFunc(handlers.HandlerNotFound)

	// Handle Method Not Allowed
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.HandlerMethodNotAllowed)

	// subrouter so it can be used a version previously to any resource
	path := router.PathPrefix(data.apiVersion).Subrouter()

	// request handler resource
	path.Use(handlers.HandlerRequestHandler)

	// log.Println(auth.AuthorizationBearerToken(http.))

	// index resource
	path.HandleFunc("/", resources.ResourceIndex).Methods("GET")

	// users resource
	path.HandleFunc("/users/{id}", resources.ResourceUsers).Methods("POST")

	// print text to let knoe the server is running
	log.Println("Listenting on Port: " + data.port)

	// start server or log error
	err := http.ListenAndServe(":"+data.port, router)

	if err != nil {
		log.Fatal("Server Start Error: " + err.Error())
	}

}
