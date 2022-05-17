package main

import (
	"challange/rest"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	//Recebe e Insere pedido no banco
	r.HandleFunc("/ordem", rest.OrdemHandler)
	http.Handle("/", router)

	log.Print("Padrozinando para porta 5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}

}
