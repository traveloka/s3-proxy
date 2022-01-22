package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/traveloka/s3-proxy/auth"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/_healthz", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "ready")
	})

	mux.HandleFunc("/_auth", auth.HandleAuthCallback)
	mux.HandleFunc("/_logout", auth.HandleLogout)
	mux.Handle("/", auth.WithAuth(http.HandlerFunc(handleProxy)))
	port := "8080"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	fmt.Printf("listening on http://127.0.0.0:%s/", port)
	http.ListenAndServe(":"+port, mux)
}
