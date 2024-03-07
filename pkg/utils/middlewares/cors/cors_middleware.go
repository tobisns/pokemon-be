package cors

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func MiddleCORS(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter,
		r *http.Request, ps httprouter.Params) {
		// Set CORS headers
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		log.Println("this was called")

		next(w, r, ps)
	}
}
