package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type MyMiddleWare struct {
	hanler http.Handler
}

func (mw *MyMiddleWare) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	mw.hanler.ServeHTTP(w, req)
	log.Println("time used: ", time.Since(start))
}

type PoliteServer struct {
}

func (ms *PoliteServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome! Thanks for visiting!\n")
}

func main() {
	ps := &PoliteServer{}
	mw := &MyMiddleWare{ps}
	log.Fatal(http.ListenAndServe(":8090", mw))
}
