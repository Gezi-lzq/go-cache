package lru

import "container/list"

type Cache struct {
	// 允许最大使用内存
	maxBytes int64
	// 当前已使用内存
	nbytes int64
	// 双向链表
	ll *list.List

	cache map[string]*list.Element

	// OnEvicted 是某条记录被移除时的回调函数，可以为 nil
	onEvicted func(key string, value Value)
}

// 键值对 entry 是双向链表节点的数据类型
type entry struct {
	// 在链表中仍保存每个值对应的 key 的好处在于
	// 淘汰队首节点时，需要用 key 从字典中删除对应的映射
	key   string
	value Value
}

// 为了通用性，我们允许值是实现了 Value 接口的任意类型
// 该接口只包含了一个方法 Len() int
// 用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// 实例化 Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 第一步：从字典中找到对应节点的双向链表的节点
	if ele, ok := c.cache[key]; ok {
		// 将该节点移动到队尾
		c.ll.MoveToFront(ele)
		// 返回查找到的值
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除 缓存淘汰，即移除最近最少访问的节点(队首)
func (c *Cache) RemoveOldest() {
	// 取出队首节点,从链表中删除
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		// 从字典c.cache删除该节点的映射关系
		delete(c.cache, kv.key)
		// 更新当前所用的内存 c.nbytes
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 若回调函数OnEvited不为nil,则调用回调函数
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// 新增/修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 如果键存在，则更新对应节点的值，并将该节点移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 不存在则是新增场景
		// 首先队尾添加新节点&entry{key,value}
		ele := c.ll.PushFront(&entry{key, value})
		// 字典中添加key和节点的映射关系
		c.cache[key] = ele
		// 更新从c.nbytes
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 如果超过了设定的最大值c.maxBytes,则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
