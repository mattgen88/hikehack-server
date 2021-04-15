package main

import (
	"log"
	"net"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// registers with database/sql
	_ "github.com/lib/pq"
	"github.com/mattgen88/hikehack-server/handlers"
	"github.com/mattgen88/hikehack-server/middleware"
	"github.com/mattgen88/hikehack-server/models"
	"github.com/mattgen88/hikehack-server/util"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()

	// Gather configuration
	viper.BindEnv("database_url")
	dsn := viper.GetString("database_url")

	viper.BindEnv("port")
	viper.SetDefault("port", "8088")
	port := viper.GetString("port")

	viper.BindEnv("host")
	viper.SetDefault("host", "0.0.0.0")
	host := viper.GetString("host")

	viper.BindEnv("jwtKey")
	jwtKey := viper.GetString("jwtKey")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.Trails{})

	dbh, _ := db.DB()
	defer dbh.Close()

	r := mux.NewRouter()

	h := handlers.New(r, jwtKey, db)
	cors := Gorilla.CORS(
		Gorilla.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		Gorilla.AllowedOrigins([]string{"*"}),
		Gorilla.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
		Gorilla.AllowCredentials())

	r.Handle(
		"/",
		Gorilla.MethodHandler{
			"GET": util.ContentType(h.GetRoot, "application/hal+json"),
		}).
		Name("root")
	r.Handle(
		"/trails",
		Gorilla.MethodHandler{
			"GET":  util.ContentType(h.GetTrails, "application/hal+json"),
			"POST": middleware.AuthMiddleware(util.ContentType(h.CreateTrail, "application/hal+json"), jwtKey, db),
		}).
		Name("trails")
	r.Handle(
		"/trails/{id}",
		Gorilla.MethodHandler{
			"GET": util.ContentType(h.GetTrail, "application/gpx+xml"),
		}).
		Name("Trail")
	r.Handle(
		"/auth",
		Gorilla.MethodHandler{
			"POST": util.ContentType(h.Auth, "application/hal+json"),
		}).
		Name("Auth")
	r.Handle(
		"/auth/refresh",
		Gorilla.MethodHandler{
			"POST": util.ContentType(h.AuthRefresh, "application/hal+json"),
		}).
		Name("AuthRefresh")
	r.Handle(
		"/auth/register",
		Gorilla.MethodHandler{
			"POST": util.ContentType(h.Register, "application/hal+json"),
		}).
		Name("Register")

	r.NotFoundHandler = http.HandlerFunc(handlers.Error)

	log.Fatal(
		http.ListenAndServe(
			net.JoinHostPort(host, port),
			Gorilla.LoggingHandler(
				os.Stdout,
				cors(r))))
}
