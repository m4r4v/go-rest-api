package auth

import (
	"log"
	"strings"
)

func AuthorizationBearerToken(t string) bool {

	token := strings.Split(t, "Bearer")

	if len(token) != 2 {
		log.Fatal("Error in Bearer Token")
	}

	if len(strings.TrimSpace(token[1])) != 17 {
		return false
	}

	return true

}
