package models

import "gorm.io/gorm"

type Trails struct {
	gorm.Model
	Owner   *User
	OwnerID int
	Name    string
	Title   string
	GPX     string
}
