package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/lcafe/a_simple_webserver/internal/config"
)

// ProxyHandler cria um handler HTTP que encaminha requisições para os backends configurados
func ProxyHandler(cfg *config.Config) http.Handler {
	// Normaliza o prefixo para garantir que comece com / e termine com /
	if !strings.HasPrefix(cfg.ProxyPrefix, "/") {
		cfg.ProxyPrefix = "/" + cfg.ProxyPrefix
	}
	if !strings.HasSuffix(cfg.ProxyPrefix, "/") {
		cfg.ProxyPrefix += "/"
	}

	// Cria um mapa para armazenar os proxies
	proxyMap := make(map[string]*httputil.ReverseProxy)

	// Configura os proxies para cada VirtualHost
	for path, targetURL := range cfg.VirtualHosts {
		// Remove a barra inicial do path se existir
		path = strings.TrimPrefix(path, "/")

		// Cria o caminho completo com o prefixo
		fullPath := path

		// Faz o parsing da URL de destino
		target := parseURL(targetURL)

		// Cria o proxy reverso
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Configura o proxy para modificar o caminho da requisição
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			// Extrai o caminho relativo após o prefixo
			pathPrefix := cfg.ProxyPrefix + fullPath
			req.URL.Path = strings.TrimPrefix(req.URL.Path, pathPrefix)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
		}

		// Adiciona ao mapa de proxies
		proxyMap[fullPath] = proxy
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove o prefixo da URL para encontrar o caminho relativo
		path := strings.TrimPrefix(r.URL.Path, cfg.ProxyPrefix)

		// Encontra o proxy mais longo que corresponde ao início do caminho
		var matchedPath string
		var matchedProxy *httputil.ReverseProxy

		for proxyPath, proxy := range proxyMap {
			// Verifica se o caminho da requisição começa com o caminho do proxy
			if strings.HasPrefix(path, proxyPath) {
				// Se encontrarmos um caminho mais longo, usamos ele
				if len(proxyPath) > len(matchedPath) {
					matchedPath = proxyPath
					matchedProxy = proxy
				}
			}
		}

		// Se encontramos um proxy correspondente, encaminhamos a requisição
		if matchedProxy != nil {
			matchedProxy.ServeHTTP(w, r)
			return
		}

		// Se nenhum proxy corresponder, retorna 404
		http.NotFound(w, r)
	})
}

func parseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("URL inválida: %v", err)
	}
	return u
}
