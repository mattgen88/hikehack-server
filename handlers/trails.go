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

	trails_files, err := ioutil.ReadDir("trails")
	if err != nil {
		panic(err)
	}
	var trails []string
	for _, trail := range trails_files {
		name := trail.Name()[:len(trail.Name())-4]
		root.AddLink(name, &haljson.Link{Href: "/trails/" + name})
		trails = append(trails, name)
	}
	root.Data["trails"] = trails

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}