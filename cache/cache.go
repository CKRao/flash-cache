package cache

import "log"

//Cache 缓存定义的接口
type Cache interface {
	Set(string, []byte) error   //Set方法定义，设置缓存
	Get(string) ([]byte, error) //Get方法定义，获取缓存
	Del(string) error           //Del方法定义，删除缓存
	GetStat() Stat              //GetStat方法定义，获取当前缓存的状态
}

//Stat 缓存状态结构体
type Stat struct {
	Count     int64 //缓存键值对的总数量
	KeySize   int64 //key的总字节数
	ValueSize int64 //value的总字节数
}

//增加
func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

//删除
func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}

//New 新建缓存的方法
func New(typ string) Cache {
	var c Cache

	if typ == InMemory {
		c = newInMemoryCache()
	}

	if c == nil {
		panic("unknown cache type " + typ)
	}

	log.Println(typ, " ready to serve")

	return c
}
