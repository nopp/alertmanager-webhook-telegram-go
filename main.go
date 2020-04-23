package main

import (
	"alertmanager-webhook-telegram-go/alert"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/alert", alert.ToTelegram).Methods("POST")

	log.Fatal(http.ListenAndServe("0.0.0.0:9229", router))
}
