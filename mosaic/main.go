package main

import (
	"fmt"
	"mosaic/sync"
	"mosaic/tilesdb"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", sync.Upload)
	mux.HandleFunc("/mosaic", sync.Mosaic)
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	tilesdb.TILESDB = tilesdb.TilesDB()
	fmt.Println("图片马赛克应用服务器已启动...")
	server.ListenAndServe()
}
