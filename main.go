package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	Gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// registers with database/sql
	_ "github.com/lib/pq"
	"github.com/mattgen88/hikehack-server/handlers"
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

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := mux.NewRouter()

	h := handlers.New(r, db)

	r.HandleFunc("/", h.RootHandler).Name("root")

	r.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	log.Fatal(http.ListenAndServe(net.JoinHostPort(host, port), util.ContentType(Gorilla.LoggingHandler(os.Stdout, Gorilla.CORS()(r)), "application/hal+json")))
}
