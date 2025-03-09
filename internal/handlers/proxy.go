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
	normalizeProxyPrefix(cfg)
	proxyMap := buildProxyMap(cfg)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matchedProxy := matchProxy(cfg.ProxyPrefix, proxyMap, r.URL.Path)
		if matchedProxy != nil {
			matchedProxy.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

func normalizeProxyPrefix(cfg *config.Config) {
	if !strings.HasPrefix(cfg.ProxyPrefix, "/") {
		cfg.ProxyPrefix = "/" + cfg.ProxyPrefix
	}
	if !strings.HasSuffix(cfg.ProxyPrefix, "/") {
		cfg.ProxyPrefix += "/"
	}
}

func buildProxyMap(cfg *config.Config) map[string]*httputil.ReverseProxy {
	proxyMap := make(map[string]*httputil.ReverseProxy)

	for path, targetURL := range cfg.VirtualHosts {
		normalizedPath := strings.TrimPrefix(path, "/")
		target := parseURL(targetURL)
		proxy := httputil.NewSingleHostReverseProxy(target)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			adjustRequestPath(req, cfg.ProxyPrefix, normalizedPath)
		}
		proxyMap[normalizedPath] = proxy
	}

	return proxyMap
}

func adjustRequestPath(req *http.Request, proxyPrefix, path string) {
	pathPrefix := proxyPrefix + path
	req.URL.Path = strings.TrimPrefix(req.URL.Path, pathPrefix)
	if !strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = "/" + req.URL.Path
	}
}

func matchProxy(proxyPrefix string, proxyMap map[string]*httputil.ReverseProxy, requestPath string) *httputil.ReverseProxy {
	trimmedPath := strings.TrimPrefix(requestPath, proxyPrefix)
	var matchedPath string
	var matchedProxy *httputil.ReverseProxy

	for path, proxy := range proxyMap {
		if strings.HasPrefix(trimmedPath, path) && len(path) > len(matchedPath) {
			matchedPath = path
			matchedProxy = proxy
		}
	}
	return matchedProxy
}

func parseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("URL inválida: %v", err)
	}
	return u
}
