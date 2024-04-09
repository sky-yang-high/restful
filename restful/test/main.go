package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 简化版的中间件
func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println("time used:", time.Since(start))
	})
}

type PoliteServer struct {
}

func (ms *PoliteServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome! Thanks for visiting!\n")
}

func main() {
	ps := &PoliteServer{}
	//在这里，我们并没有使用到任何路由，但是这的确跑起来了
	//如果不理解这部分，可以再深入了解了解go的http包
	log.Fatal(http.ListenAndServe(":8090", loggingMiddleWare(ps)))
}
