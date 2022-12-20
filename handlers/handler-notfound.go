package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func HandlerNotFound(w http.ResponseWriter, r *http.Request) {

	httpStatus := http.StatusNotFound

	var response *DefaultResponse = &DefaultResponse{
		Status:  strconv.Itoa(httpStatus),
		Message: "Error 404, your request was not found",
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.Fatal("Error 404: " + err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(jsonResponse)

}
