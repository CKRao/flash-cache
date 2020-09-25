package main

import (
	"flash-cache/cache"
	"flash-cache/http"
)

//flash-cache is a fast as flash distribution cache
//powered by clarkRao

func main() {
	c := cache.New(cache.InMemory)
	http.New(c).Listen()
}
