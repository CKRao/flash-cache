package main

import (
	"flash-cache/cache"
	"flash-cache/http"
	"flash-cache/tcp"
)

//flash-cache is a fast as flash distribution cache
//powered by clarkRao

func main() {
	c := cache.New(cache.InMemory)

	//启动tcp的server
	go tcp.New(c).Listen()

	//http server 开始监听
	http.New(c).Listen()
}
