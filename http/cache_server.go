package http

import (
	"flash-cache/cache"
	"net/http"
)

const (
	defaultPort = ":12345" //默认端口
)

//缓存服务器结构体
type Server struct {
	cache.Cache //缓存接口
}

//Listen 缓存服务监听
func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())

	http.ListenAndServe(defaultPort, nil)
}

//新建Server
func New(c cache.Cache) *Server {
	return &Server{c}
}
