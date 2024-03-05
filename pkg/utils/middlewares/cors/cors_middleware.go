package cors

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func MiddleCORS(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter,
		r *http.Request, ps httprouter.Params) {
		// Set CORS headers
		w.Header().Add("Access-Control-Allow-Origin", "*")

		// Second, we handle the OPTIONS problem
		if r.Method != "OPTIONS" {

			// Call the next handler
			next(w, r, ps)

		} else {

			// Return HTTP 200 OK for OPTIONS requests
			w.WriteHeader(http.StatusOK)
		}

	}
}
