package resources

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	handlers "github.com/jmr-repo/go-rest-api/handlers"
)

func ResourceIndex(w http.ResponseWriter, r *http.Request) {

	httpStatus := http.StatusOK

	var response *handlers.DefaultResponse = &handlers.DefaultResponse{
		Status:  strconv.Itoa(httpStatus),
		Message: "hello world!",
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.Fatal("jsonResponse Error: " + err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(jsonResponse)

}
