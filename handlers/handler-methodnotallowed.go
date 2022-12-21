package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmr-repo/go-rest-api/interfaces"
)

func HandlerMethodNotAllowed(w http.ResponseWriter, r *http.Request) {

	httpStatus := http.StatusMethodNotAllowed

	response = &interfaces.IDefaultResponse{
		Status:  httpStatus,
		Message: "Error 405, your request method is not allowed",
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.Fatal("Error 405: " + err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(jsonResponse)

}
