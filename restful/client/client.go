package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:8080/secret", "server address")
	user := flag.String("user", "admin", "username")
	password := flag.String("password", "123456", "password")
	certFile := flag.String("cert", "../tls/cert.pem", "certificate file path")
	flag.Parse()

	//load the certificate
	//NOTE: the certificate is used to verify the server's identity
	//not for client authentication
	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	certpool := x509.NewCertPool()
	if ok := certpool.AppendCertsFromPEM(cert); !ok {
		log.Fatal("failed to parse certificate")
	}

	//create a new http client with the certificate pool
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certpool,
			},
		},
	}

	//do request
	req, err := http.NewRequest("GET", "https://"+*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	//automatically encode the username and password in base64 format
	req.SetBasicAuth(*user, *password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTTP Status:", resp.Status)
	fmt.Println("Response Body:", string(msg))
}
