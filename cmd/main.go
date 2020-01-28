package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fernandoguedes/rgb-soccer"
)

const port = ":8081"

func main() {
	http.HandleFunc("/", rigabe.Rigabe)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(port, nil))
}
