package net

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var Addr = ":8081"

func HttpTest() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 休眠 5s
		time.Sleep(6 * time.Second)
		fmt.Println("Received request from:", r.RemoteAddr)
		w.Write([]byte("Hello, World\n"))
	})

	server := &http.Server{
		Addr:         Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	log.Println("Server started on", Addr)
	log.Fatal(server.ListenAndServe())

	// http.ListenAndServe(Addr, mux)
}
