package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetTrail handles requests for the root of the API
func (h *Handler) GetTrail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	trail, err := ioutil.ReadFile("trails/" + name + ".gpx")
	if err != nil {
		// error
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(trail)
}
