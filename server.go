package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/v1/readiness", readiness)

	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: router,
	}

	go func() {
		log.Println("server listening on", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdownChannel
	log.Println("signal:", sig)

	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		server.Close()
	}

}

func readiness(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-REQUEST-ID")
	log.Println("start", requestID)
	defer log.Println("done", requestID)

	time.Sleep(5 * time.Second)

	response := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		panic(err)
	}
}
