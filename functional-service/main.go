package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	nat "github.com/ravjotsingh9/nats-streaming-functional-service/nats"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/logs", getLogsHandler).
		Methods("GET")
	return
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	//connect to db
	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	db, err := sql.Open("postgres", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	//connect to nat
	/*	es, err = nats.Connect(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			return
		}
	*/

	es, err := nat.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Println(err)
		return
	}
	err = es.OnLogsCreated(func(m nat.LogsCreatedMessage) {
		fmt.Println("OnLogsCreated " + m.LogContent)
		db.Exec("INSERT INTO distributedlogs(id, logContent, created_at) VALUES($1, $2, $3)", m.ID, m.LogContent, m.CreatedAt)

		row, err := db.Exec("SELECT * FROM distributedlogs")
		fmt.Println(row)
		fmt.Println(err)
	})
	if err != nil {
		log.Println(err)
		return
	}
	// Run HTTP server

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

}
