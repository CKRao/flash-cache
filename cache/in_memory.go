package cache

import "sync"

//内存缓存结构体
type inMemoryCache struct {
	c     map[string][]byte //存储键值对的map
	mutex sync.RWMutex      //读写锁
	Stat                    //缓存状态
}

//Set 设置缓存
func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	tmp, exist := c.c[k]

	//如果已经存在，则先删除缓存状态
	if exist {
		c.del(k, tmp)
	}

	//设置val
	c.c[k] = v
	//重新更新缓存状态
	c.add(k, v)

	return nil
}

//Get 获取缓存的方法
func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k], nil
}

//Del 删除缓存的方法
func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v, exist := c.c[k]

	//如果缓存存在，则删除
	if exist {
		delete(c.c, k)
		c.del(k, v)
	}

	return nil
}

//GetStat 获取缓存的状态
func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

//新建内存缓存
func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		c:     make(map[string][]byte),
		mutex: sync.RWMutex{},
		Stat:  Stat{},
	}
}
