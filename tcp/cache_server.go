package tcp

import (
	"bufio"
	"flash-cache/cache"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	defaultPort = ":12346" //默认端口
)

//TCP Server 结构体
type Server struct {
	cache.Cache //缓存
}

//Listen 实现监听方法
func (s *Server) Listen() {
	l, e := net.Listen("tcp", defaultPort)

	if e != nil {
		panic(e)
	}

	for {
		c, e := l.Accept()

		if e != nil {
			panic(e)
		}

		go s.process(c)
	}
}

//读取key的方法
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, e := readLen(r)

	if e != nil {
		return "", e
	}

	k := make([]byte, klen)

	_, e = io.ReadFull(r, k)

	if e != nil {
		return "", e
	}

	return string(k), nil
}

//读取key和value
func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	kLen, e := readLen(r)

	if e != nil {
		return "", nil, e
	}

	vLen, e := readLen(r)

	if e != nil {
		return "", nil, e
	}

	k := make([]byte, kLen)

	_, e = io.ReadFull(r, k)

	if e != nil {
		return "", nil, e
	}

	v := make([]byte, vLen)

	_, e = io.ReadFull(r, v)

	if e != nil {
		return "", nil, e
	}

	return string(k), v, nil
}

//读取长度的函数
func readLen(r *bufio.Reader) (int, error) {
	tmp, e := r.ReadString(' ')

	if e != nil {
		return 0, e
	}

	l, e := strconv.Atoi(strings.TrimSpace(tmp))

	if e != nil {
		return 0, e
	}

	return l, nil
}

//发送响应
func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()

		tmp := fmt.Sprintf("-%d %s", len(errString), errString)

		_, e := conn.Write([]byte(tmp))

		return e
	}

	vLen := fmt.Sprintf("%d ", len(value))

	_, e := conn.Write(append([]byte(vLen), value...))

	return e
}

//获取值
func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)

	if e != nil {
		return e
	}

	v, e := s.Get(k)

	return sendResponse(v, e, conn)
}

//设置键值对
func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, e := s.readKeyAndValue(r)

	if e != nil {
		return e
	}

	return sendResponse(nil, s.Set(k, v), conn)
}

//删除键值对
func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)

	if e != nil {
		return e
	}

	return sendResponse(nil, s.Del(k), conn)
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	for {
		op, e := r.ReadByte()

		if e != nil {

			if e != io.EOF {
				log.Println("close connection due to error: ", e)
			}

			return
		}

		switch op {
		case Set:
			e = s.set(conn, r)
		case Get:
			e = s.get(conn, r)
		case Del:
			e = s.del(conn, r)
		default:
			log.Println("close connection due to invalid operation: ", op)
		}

		if e != nil {
			log.Println("close connection due to error: ", e)
		}
	}
}

//New 新建Server的方法
func New(c cache.Cache) *Server {
	return &Server{c}
}
