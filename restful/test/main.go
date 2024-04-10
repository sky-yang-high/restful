package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func doPanic(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "panic")
		}
	}()
	//panic will close the connection and shutdown the goroutine,
	//which results in the client didn't receive any response.
	//(but luckily, the server isn't impacted, since the routine is seperated from the main routine)
	//so here we just log the panic but don't do an panic itself.
	//log.Println("panic opps")
	panic("oops")
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/panic", doPanic)

	http.ListenAndServe(":8090", nil)
}
