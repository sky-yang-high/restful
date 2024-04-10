package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8099", "server address")
	certFile := flag.String("certfile", "./cert.pem", "certificate PEM file")
	privateKey := flag.String("privatekey", "./key.pem", "private key PEM file")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("Hello, TLS!"))
	})

	srv := http.Server{
		Addr:    *addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		}}

	log.Println("Starting server on", *addr)
	log.Fatal(srv.ListenAndServeTLS(*certFile, *privateKey))
}
