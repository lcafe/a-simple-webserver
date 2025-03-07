package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lcafe/a_simple_webserver/internal/config"
	"github.com/lcafe/a_simple_webserver/internal/handlers"
)

func main() {

	// Carrega a configuração
	cfg, err := config.LoadConfig("internal/config/config.json")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      mux,
		ErrorLog:     log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	mux.Handle("/apps/", handlers.ProxyHandler(cfg))
	mux.Handle("/files/", handlers.FileHandler(cfg))

	error := srv.ListenAndServe()
	log.Fatal(error)

}
