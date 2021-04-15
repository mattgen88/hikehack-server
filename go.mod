module github.com/mattgen88/hikehack-server

// +heroku goVersion 1.16
go 1.16

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.0
	github.com/mattgen88/haljson v1.0.2
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.7
)
