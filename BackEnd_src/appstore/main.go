package main

import (
	"fmt"
	"log"
	"net/http"

	"appstore/backend"
	"appstore/handler"
)
func main() {
    fmt.Println("started-service")
    backend.InitElasticsearchBackend()//先初始化
	backend.InitGCSBackend()
    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))//再监听端口;
}

