package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"

	nat "github.com/ravjotsingh9/nats-streaming-functional-service/nats"
	"github.com/ravjotsingh9/nats-streaming-functional-service/util"
	"github.com/segmentio/ksuid"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/", createLogHandler).
		Methods("POST")
	router.HandleFunc("/publishlog/", createLogHandler).
		Methods("POST")
	return router
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Nats
	es, err := nat.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Println(err)
		return
	}
	nat.SetNatsInterface(es)
	fmt.Println("Connected to Nats:" + cfg.NatsAddress)
	// Run HTTP server
	router := newRouter()
	if err := http.ListenAndServe(":1234", router); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println("Started http server.")
}

func createLogHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In handler")
	type response struct {
		ID string `json:"id"`
	}

	//ctx := r.Context()

	// Read parameters
	body := template.HTMLEscapeString(r.FormValue("body"))
	fmt.Println("BODY:" + body)
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	// Create meow
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create meow")
		return
	}
	obj := nat.LogsCreatedMessage{
		ID:         id.String(),
		LogContent: body,
		CreatedAt:  createdAt,
	}

	// Publish event
	if err = nat.PublishLogsCreated(obj); err != nil {
		log.Println(err)
	}

	// Return new meow
	util.ResponseOk(w, response{ID: obj.ID})
}
