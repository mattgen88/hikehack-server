package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
)

// RootHandler handles requests for the root of the API
func (h *Handler) TrailsHandler(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self("/trails")

	trails, err := ioutil.ReadDir("trails")
	if err != nil {
		panic(err)
	}
	for _, trail := range trails {
		name := trail.Name()[:len(trail.Name())-4]
		root.AddLink("trail", &haljson.Link{Href: "/trails/" + name})
		root.Data[name] = name
	}

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
