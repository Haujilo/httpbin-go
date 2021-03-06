package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	addr := "0.0.0.0:1121"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	mux := http.NewServeMux()
	route(mux)

	log.Println("Starting httpbin", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
