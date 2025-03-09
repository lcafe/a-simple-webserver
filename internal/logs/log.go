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
	logHeaders(r.Header)
	logRequestBody(r)
}

func logHeaders(headers http.Header) {
	for key, values := range headers {
		log.Printf("Header: %s = %v", key, values)
	}
}

func logRequestBody(r *http.Request) {
	if r.Body == nil {
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Erro ao ler body: %v", err)
		r.Body = io.NopCloser(strings.NewReader(""))
		return
	}
	defer r.Body.Close()

	bodyStr := maskSensitiveData(string(bodyBytes))
	log.Printf("Body: %s", bodyStr)

	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
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
	bodyStr := maskSensitiveData(string(lrw.body))
	log.Printf("<<< Resposta enviada: status: %d, body: %s", lrw.StatusCode, bodyStr)
}
