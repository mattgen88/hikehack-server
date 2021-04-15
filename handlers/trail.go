package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattgen88/haljson"
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

// CreateTrail handles requests to create a trail
func (h *Handler) CreateTrail(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self(r.URL.Path)

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
