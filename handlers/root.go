package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
)

// Root handles requests for the root of the API
func (h *Handler) GetRoot(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self("/")

	title := "Trail list"
	root.AddLink("nav", &haljson.Link{Href: "/trails", Title: &title})

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
