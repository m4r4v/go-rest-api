package resources

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/m4r4v/go-rest-api/auth"
	interfaces "github.com/m4r4v/go-rest-api/interfaces"
)

type PostData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var responseUsers *interfaces.IDefaultResponse

func ResourceUsers(w http.ResponseWriter, r *http.Request) {

	// check if user is authorized or authenticated
	if !auth.AuthorizationBearerToken(r.Header.Get("Authorization")) {

		responseUsers = &interfaces.IDefaultResponse{
			Status:  http.StatusForbidden,
			Message: "Error 403, you do no have permission to access this resource",
		}

		log.Println("Index Forbidden")

	} else {

		decoder := json.NewDecoder(r.Body)

		var post PostData
		err := decoder.Decode(&post)

		if err != nil {
			panic(err)
		}

		if post.Username != "nano@gmail.com" {
			responseUsers = &interfaces.IDefaultResponse{
				Status:  http.StatusForbidden,
				Message: "Tu nombre de usuario es erroneo",
			}
		}

		responseUsers = &interfaces.IDefaultResponse{
			Status:  http.StatusOK,
			Message: "username: " + post.Username + ", password: " + post.Password,
		}

		log.Println("username: " + post.Username + ", password: " + post.Password)

	}

	jsonResponse, err := json.Marshal(responseUsers)

	if err != nil {
		log.Fatal("jsonResponse Error: " + err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseUsers.Status)
	w.Write(jsonResponse)

}
