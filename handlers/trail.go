package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattgen88/haljson"
	"github.com/mattgen88/hikehack-server/middleware"
	"github.com/mattgen88/hikehack-server/models"
)

// GetTrail handles requests for the root of the API
func (h *Handler) GetTrail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	trail := &models.Trails{}
	h.db.Where("id = ?", id).First(trail)
	if trail.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write(trail.GetGPX().Bytes())
}

// CreateTrail handles requests to create a trail
func (h *Handler) CreateTrail(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()

	root.Self(r.URL.Path)

	user_claims := (r.Context().Value(middleware.UserDataKey("user_data"))).(middleware.Claims)
	user := &models.User{}
	h.db.Where("username = ?", user_claims.Username).First(user)

	trail := models.Trails{}

	read_form, err := r.MultipartReader()
	for {
		part, err_part := read_form.NextPart()
		if err_part == io.EOF {
			break
		}
		if part.FormName() == "file" {
			// do something with files
			// Not a valid type
			if part.Header.Get("Content-Type") != "application/gpx+xml" {
				root.Data["error"] = "Only GPX is supported"
				w.WriteHeader(http.StatusBadRequest)
				json, marshalErr := json.Marshal(root)
				if marshalErr != nil {
					log.Println(marshalErr)
					return
				}
				w.Write(json)
				return
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			trail.SetGPX(buf)

		} else if part.FormName() == "name" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			trail.Name = buf.String()
		}
	}
	trail.Owner = user

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
		return
	}

	h.db.Create(&trail)

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
