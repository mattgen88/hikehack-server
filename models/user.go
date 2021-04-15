package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Realname string
	Role     string
	Email    string `gorm:"uniqueIndex"`
	Password string
}

func (u *User) Authenticate(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if err != nil {
		log.Println("mismatched hashing")
	}
	return err == nil
}
