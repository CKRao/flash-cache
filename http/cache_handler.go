package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//缓存Handler
type cacheHandler struct {
	*Server
}

//实现缓存的服务方法
func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]

	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := r.Method

	//如果是PUT方法，则说明是设置缓存的操作
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) == 0 {
			return
		}

		e := h.Set(key, b)
		if e != nil {
			log.Println(fmt.Sprintf("cache Set fail, key: %s, err: %s", key, e.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	//如果是Get方法，则说明是获取缓存的操作
	if m == http.MethodGet {
		b, e := h.Get(key)

		if e != nil {
			log.Println(fmt.Sprintf("cache Get fail, key: %s, err: %s", key, e.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(b) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write(b)
		return
	}

	//如果是Delete方法，则说明是删除缓存的操作
	if m == http.MethodDelete {
		e := h.Del(key)

		if e != nil {
			log.Println(fmt.Sprintf("cache Del fail, key: %s, err: %s", key, e.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}
