package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// RootHandler handles requests for the root of the API
func (h *Handler) TrailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	trail, err := ioutil.ReadFile("trails/" + name + ".gpx")
	if err != nil {
		// error
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/xml")
	w.Write(trail)
}
