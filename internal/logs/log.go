package logs

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func maskSensitiveData(input string) string {
	replacer := strings.NewReplacer(
		"senha", "[MASK]",
		"token", "[MASK]",
		"password", "[MASK]",
		"login", "[MASK]",
	)
	return replacer.Replace(input)
}

func LogRequest(r *http.Request) {
	log.Println(">>> Requisição recebida:")
	log.Printf("Método: %s, URL: %s", r.Method, r.URL.String())
	for k, v := range r.Header {
		log.Printf("Header: %s = %v", k, v)
	}

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			defer r.Body.Close()
			bodyStr := string(bodyBytes)
			bodyStr = maskSensitiveData(bodyStr)
			log.Printf("Body: %s", bodyStr)
			r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		} else {
			log.Printf("Erro ao ler body: %v", err)
		}
	}
}

type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	body       []byte
}

func (lrw *LogResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LogResponseWriter) Write(b []byte) (int, error) {
	lrw.body = append(lrw.body, b...)
	return lrw.ResponseWriter.Write(b)
}

func LogResponse(lrw *LogResponseWriter) {
	bodyStr := string(lrw.body)
	bodyStr = maskSensitiveData(bodyStr)
	log.Printf("<<< Resposta enviada: status: %d, body: %s", lrw.StatusCode, bodyStr)
}
