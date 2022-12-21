package resources

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/jmr-repo/go-rest-api/auth"
	interfaces "github.com/jmr-repo/go-rest-api/interfaces"
)

var response *interfaces.IDefaultResponse

func ResourceIndex(w http.ResponseWriter, r *http.Request) {

	// check if user is authorized or authenticated
	if !auth.AuthorizationBearerToken(r.Header.Get("Authorization")) {

		response = &interfaces.IDefaultResponse{
			Status:  http.StatusForbidden,
			Message: "Error 403, you do no have permission to access this resource",
		}

		log.Println("Index Forbidden")

	} else {

		response = &interfaces.IDefaultResponse{
			Status:  http.StatusOK,
			Message: "Hello world!",
		}

		log.Println("Index")

	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.Fatal("jsonResponse Error: " + err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	w.Write(jsonResponse)

}
