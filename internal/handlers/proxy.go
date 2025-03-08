package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/lcafe/a_simple_webserver/internal/config"
)

func ProxyHandler(cfg *config.Config) http.Handler {
	if !strings.HasPrefix(cfg.ProxyPrefix, "/") {
		cfg.ProxyPrefix = "/" + cfg.ProxyPrefix
	}
	if !strings.HasSuffix(cfg.ProxyPrefix, "/") {
		cfg.ProxyPrefix += "/"
	}

	proxyMap := make(map[string]*httputil.ReverseProxy)

	for path, targetURL := range cfg.VirtualHosts {
		path = strings.TrimPrefix(path, "/")

		fullPath := path

		target := parseURL(targetURL)

		proxy := httputil.NewSingleHostReverseProxy(target)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			pathPrefix := cfg.ProxyPrefix + fullPath
			req.URL.Path = strings.TrimPrefix(req.URL.Path, pathPrefix)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
		}

		proxyMap[fullPath] = proxy
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, cfg.ProxyPrefix)

		var matchedPath string
		var matchedProxy *httputil.ReverseProxy

		for proxyPath, proxy := range proxyMap {
			if strings.HasPrefix(path, proxyPath) {
				if len(proxyPath) > len(matchedPath) {
					matchedPath = proxyPath
					matchedProxy = proxy
				}
			}
		}

		if matchedProxy != nil {
			matchedProxy.ServeHTTP(w, r)
			return
		}

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
