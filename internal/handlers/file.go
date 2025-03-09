package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lcafe/a_simple_webserver/internal/config"
)

type DirEntry struct {
	Name    string
	URL     string
	Size    string
	ModTime string
}

type DirListingData struct {
	Dir       string
	ParentURL string
	Entries   []DirEntry
}

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

func buildDirListingData(directory, urlPrefix string, r *http.Request) (DirListingData, error) {
	data := DirListingData{
		Dir:     r.URL.Path,
		Entries: []DirEntry{},
	}

	currentURL := r.URL.Path
	if !strings.HasSuffix(currentURL, "/") {
		currentURL += "/"
	}

	if currentURL != urlPrefix {
		data.ParentURL = "../"
	}

	dir, err := os.Open(directory)
	if err != nil {
		return data, err
	}
	defer dir.Close()

	infos, err := dir.Readdir(-1)
	if err != nil {
		return data, err
	}

	for _, info := range infos {
		entryURL := currentURL + info.Name()
		if info.IsDir() {
			entryURL += "/"
		}
		size := "-"
		if !info.IsDir() {
			size = formatSize(info.Size())
		}
		entry := DirEntry{
			Name:    info.Name(),
			URL:     entryURL,
			Size:    size,
			ModTime: info.ModTime().Format("02/01/2006 15:04"),
		}
		data.Entries = append(data.Entries, entry)
	}

	return data, nil
}

func formatSize(size int64) string {
	return fmt.Sprintf("%d B", size)
}

func serveDirListing(w http.ResponseWriter, r *http.Request, fullPath, urlPrefix string) {
	tmpl, err := template.ParseFiles("./templates/folders.html")
	if err != nil {
		http.Error(w, "Erro ao carregar template.", http.StatusInternalServerError)
		log.Printf("Erro ao carregar template: %v\n", err)
		return
	}

	data, err := buildDirListingData(fullPath, urlPrefix, r)
	if err != nil {
		http.Error(w, "Erro ao ler diretório.", http.StatusInternalServerError)
		log.Printf("Erro ao ler diretório: %v\n", err)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erro ao renderizar a listagem.", http.StatusInternalServerError)
		log.Printf("Erro ao renderizar a listagem: %v\n", err)
	}
}

func FileHandler(cfg *config.Config) http.Handler {
	log.Println("Inicializando servidor de arquivos.")

	if !validateDirectory(cfg.FileServer) {
		log.Println("Utilizando handler padrão, pois o diretório de arquivos não foi validado.")
		return http.DefaultServeMux
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		relPath := strings.TrimPrefix(r.URL.Path, cfg.FileServerRootUrl)
		fullPath := filepath.Join(cfg.FileServer, filepath.FromSlash(relPath))

		info, err := os.Stat(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if info.IsDir() {
			indexPath := filepath.Join(fullPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
				return
			}
			serveDirListing(w, r, fullPath, cfg.FileServerRootUrl)
		} else {
			http.ServeFile(w, r, fullPath)
		}
	})
}
