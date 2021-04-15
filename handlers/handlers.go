package handlers

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Handler provides various http handlers
type Handler struct {
	r      *mux.Router
	db     *gorm.DB
	jwtKey string
}

// New returns a configured handler struct
func New(r *mux.Router, jwtKey string, db *gorm.DB) *Handler {
	return &Handler{r, db, jwtKey}
}
