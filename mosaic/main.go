package main

import (
	"fmt"
	"mosaic/conc"
	"mosaic/noconc"
	"mosaic/tilesdb"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", noconc.Upload)
	//这里有两个版本的Mosaic，因为不想改html文件，所以就只好改这里了
	//mux.HandleFunc("/mosaic", noconc.Mosaic)
	mux.HandleFunc("/mosaic", conc.Mosaic)
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	tilesdb.TILESDB = tilesdb.TilesDB()
	fmt.Println("图片马赛克应用服务器已启动...")
	server.ListenAndServe()
}
