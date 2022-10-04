package geecache

import (
	"fmt"
	"log"
	"sync"
)

// 定义接口Getter与回调函数Get(key string)([]byte,error)
type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义函数类型GetterFunc 并实现Getter接口的Get方法
type GetterFunc func(key string) ([]byte, error)

// 函数类型实现某一个接口,称之为接口型函数
// 方便使用者在调用时既能够传入函数作为参数
// 也能够传入实现该接口的结构体作为参数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group可以被认为一个缓存的命名空间
type Group struct {
	name string
	// 缓存未命中时获取源数据的回调
	getter Getter
	// 并发缓存
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	// 由于不涉及任何冲突变量的写操作,故可以使用只读锁RLock
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is requestd")
	}
	// 从mainCache中查找缓存,如果存在则返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 缓存不存在则调用load方法
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 分布式场景下会电影getFromPeer从其他节点获取
	// 此处调用getLocally
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用用户回调函数 g.getter.Get() 获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	// 并将源数据添加到缓存mainCache中
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
