package handlers

import (
	"net/http"

	"github.com/jmr-repo/go-rest-api/interfaces"
)

var response *interfaces.IDefaultResponse

func HandlerRequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

	})
}
