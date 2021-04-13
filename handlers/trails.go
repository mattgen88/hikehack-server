package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/mattgen88/haljson"
)

// RootHandler handles requests for the root of the API
func (h *Handler) TrailsHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self("/trails")

	trails, err := os.ReadDir("trails")
	if err != nil {
		panic(err)
	}
	for _, trail := range trails {
		root.Data[trail.Name()] = trail.Name()
	}

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
