package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lcafe/a_simple_webserver/internal/config"
)

func validateDirectory(directory string) bool {
	if directory == "" {
		log.Println("Campo 'file_server' não está configurado.")
		return false
	}

	info, err := os.Stat(directory)
	if err != nil {
		log.Printf("Erro ao acessar diretório '%s': %v\n", directory, err)
		return false
	}

	if !info.IsDir() {
		log.Printf("O caminho '%s' não é um diretório.\n", directory)
		return false
	}

	return true
}

func hasIndexFile(directory string) bool {
	indexPath := filepath.Join(directory, "index.html")
	if info, err := os.Stat(indexPath); err == nil && !info.IsDir() {
		return true
	}
	return false
}

func FileHandler(cfg *config.Config) http.Handler {
	log.Println("Inicializando servidor de arquivos.")

	if !validateDirectory(cfg.FileServer) {
		log.Println("Utilizando handler padrão, pois o diretório de arquivos não foi validado.")
		return http.DefaultServeMux
	}

	if hasIndexFile(cfg.FileServer) {
		log.Printf("Arquivo 'index.html' encontrado no diretório '%s'.\n", cfg.FileServer)
	} else {
		log.Printf("Nenhum 'index.html' encontrado no diretório '%s'. Será exibida a listagem do diretório.\n", cfg.FileServer)
	}

	handler := http.StripPrefix(cfg.FileServerRootUrl,
		http.FileServer(http.Dir(cfg.FileServer)))
	log.Printf("Servidor de arquivos configurado com root URL '%s' e diretório '%s'.\n",
		cfg.FileServerRootUrl, cfg.FileServer)

	return handler
}
