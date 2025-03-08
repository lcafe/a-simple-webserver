package handlers

import (
	"log"
	"net/http"

	"github.com/lcafe/a_simple_webserver/internal/config"
)

func FileHandler(cfg *config.Config) http.Handler {
	log.Println("Servidor de arquivos em execução.")

	if cfg.FileServer == "" {
		return http.DefaultServeMux
	}

	if _, err := http.Get(cfg.FileServer + "/index.html"); err == nil {
		return http.DefaultServeMux
	}

	return http.StripPrefix(cfg.FileServerRootUrl, http.FileServer(http.Dir(cfg.FileServer)))
}
