package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>欢迎使用稳定安全的k8s系统！</h1>")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":80", nil)
}
