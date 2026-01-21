package middleware

import (
	"bytes"
	"net/http"
)

var IdemKeyMap map[string]CachedResponse

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body bytes.Buffer
}

type CachedResponse struct {
	StatusCode int
	Body bytes.Buffer
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		statusCode: http.StatusOK,
	}
}

func (r *responseRecorder) WriterHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IdemMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip for options to declare allowed headers
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
	
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
		    http.Error(w, "Missing Idempotency-Key", http.StatusBadRequest)
		    return
		}
		
		if resp, ok := IdemKeyMap[key]; ok {
			w.WriteHeader(resp.StatusCode)
			w.Write(resp.Body.Bytes())
			return
		}
		
		recorder := newResponseRecorder(w)
		next.ServeHTTP(recorder, r)
		
		IdemKeyMap[key] = CachedResponse{
			StatusCode: recorder.statusCode,
			Body: recorder.body,
		}
    })
}
