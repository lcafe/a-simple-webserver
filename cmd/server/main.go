package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lcafe/a_simple_webserver/internal/config"
	"github.com/lcafe/a_simple_webserver/internal/handlers"
	"github.com/lcafe/a_simple_webserver/internal/middleware"
)

func main() {
	cfg, err := config.GetConfig("internal/config/config.json")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	log.Println("Servidor iniciado.")

	mux := createMux(cfg)

	srv := createServer(cfg, mux)

	go func() {
		log.Printf("Servidor rodando na porta: %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	waitForShutdown(srv, 15*time.Second)
}

func createMux(cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/{$}", middleware.LogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})))

	mux.Handle(cfg.ProxyPrefix, middleware.LogMiddleware(handlers.ProxyHandler(cfg)))

	mux.Handle(cfg.FileServerRootUrl, middleware.LogMiddleware(handlers.FileHandler(cfg)))

	return mux
}

func createServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      handler,
		ErrorLog:     log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func waitForShutdown(srv *http.Server, timeout time.Duration) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGTERM, syscall.SIGINT)

	sig := <-shutdownCh
	log.Printf("Recebido sinal %v. Iniciando shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro durante shutdown: %v", err)
	}

	log.Println("Shutdown concluído")
}
