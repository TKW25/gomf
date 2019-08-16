package main

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] STARTING REQUEST for %s\n", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
		log.Printf("[%s] COMPLETED REQUEST for %s\n", r.Method, r.URL.String())
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Catching panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})

	return http.HandlerFunc(fn)
}
