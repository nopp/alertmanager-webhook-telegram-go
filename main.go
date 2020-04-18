package main

import (
	"log"
	"net/http"

	"./alert"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/alert", alert.ToTelegram).Methods("POST")

	log.Fatal(http.ListenAndServe("0.0.0.0:8686", router))
}
