package main

import (
	"FetchChallenge/router"
	"FetchChallenge/store"
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	pHost := flag.String("h", "localhost", "host name")
	pPort := flag.Int("p", 8080, "port number")
	flag.Parse()

	store.CreateReceiptStore()
	r := router.CreateRouter()
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(*pHost+":"+strconv.Itoa(*pPort), nil))
}
