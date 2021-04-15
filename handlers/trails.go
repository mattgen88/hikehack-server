package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mattgen88/haljson"
	"github.com/mattgen88/hikehack-server/models"
)

// GetTrails handles retrieving a list of trails
func (h *Handler) GetTrails(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self(r.URL.Path)

	trails := []models.Trails{}
	h.db.Select("id", "name").Find(&trails)
	root.Data["trails"] = []struct{}{}

	var trail_list = []struct {
		Name string
		ID   uint
	}{}
	for _, trail := range trails {
		root.AddLink(fmt.Sprintf("%d", trail.ID), &haljson.Link{Href: fmt.Sprintf("/trails/%d", trail.ID)})
		trail_list = append(trail_list, struct {
			Name string
			ID   uint
		}{Name: trail.Name, ID: trail.ID})
	}
	root.Data["trails"] = trail_list

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
