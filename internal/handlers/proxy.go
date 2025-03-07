package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/lcafe/a_simple_webserver/internal/config"
)

// pathBasedHandler implementa o path-based routing.
type pathBasedHandler struct {
	cfg      *config.Config
	fallback http.Handler
}

// Uma nova função que serve os VirtualHosts definidos no arquivo de configuração.
func ProxyHandler(cfg *config.Config) http.Handler {
	// Path-based routing (para backends via VirtualHosts)
	baseHandler := &pathBasedHandler{
		cfg: cfg,
		fallback: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Caminho não cadastrado. Verifique sua configuração."))
		}),
	}

	// Para cada VirtualHost, cria um proxy reverso.
	for prefix, backend := range cfg.VirtualHosts {
		app_prefix := "/apps" + prefix
		proxy := http.StripPrefix(app_prefix, httputil.NewSingleHostReverseProxy(parseURL(backend)))
		baseHandler.fallback = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("VirtualHost: encaminhando %s para %s", app_prefix, backend)
			proxy.ServeHTTP(w, r)
		})
	}
	return baseHandler.fallback
}

func parseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("URL inválida: %v", err)
	}
	return u
}
