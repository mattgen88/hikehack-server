package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mattgen88/haljson"
	"github.com/mattgen88/hikehack-server/middleware"
	"github.com/mattgen88/hikehack-server/models"
)

// Auth handles request to authenticate and will issue a JWT
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)

	if r.Method != http.MethodPost {
		root.Data["error"] = "Please POST credentials."
		root.Data["required_fields"] = []string{"username", "password"}

		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(json)
		return
	}

	if r.FormValue("username") == "" || r.FormValue("password") == "" {
		root.Data["required_fields"] = []string{"username", "password"}
		root.Data["error"] = "Missing required fields."
		root.Data["result"] = false

		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
		return
	}

	user := &models.User{}

	h.db.Where("username = ?", r.FormValue("username")).First(user)
	if user.ID == 0 {
		// Not found
		root.Data["error"] = "Unable to authenticate. Check that credentials are correct."
		root.Data["result"] = false

		w.WriteHeader(http.StatusForbidden)
		json, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			log.Println(marshalErr)
			return
		}
		w.Write(json)
		return
	}

	now := time.Now()

	accessExpires := now.Add(time.Minute * 5)
	refreshExpires := now.Add(time.Hour * 24)

	// Create the Claims
	accessClaims := middleware.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpires.Unix(),
			Issuer:    "test",
		},
	}

	refreshClaims := middleware.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpires.Unix(),
			Issuer:    "test",
		},
	}

	accessCookie, accessErr := middleware.CreateJwt("access.jwt", accessExpires, &accessClaims, h.jwtKey)
	refreshCookie, refreshErr := middleware.CreateJwt("refresh.jwt", refreshExpires, &refreshClaims, h.jwtKey)

	if accessErr != nil {
		root.Data["err"] = fmt.Sprintf("%s", accessErr)
		root.Data["result"] = false
	} else if refreshErr != nil {
		root.Data["err"] = fmt.Sprintf("%s", refreshErr)
		root.Data["result"] = false
	} else {
		root.Data["result"] = true
		root.Data["access_expires"] = accessExpires.Unix()
		root.Data["refresh_expires"] = refreshExpires.Unix()
		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)
		root.AddLink("refresh", &haljson.Link{Href: "/refresh"})
	}

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}

// AuthRefresh returns result = true if AUTHed.
func (h *Handler) AuthRefresh(w http.ResponseWriter, r *http.Request) {
	root := haljson.NewResource()
	root.Self(r.URL.Path)
	root.Data["result"] = true

	json, err := json.Marshal(root)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(json)
}
