package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 定义了函数类型Hash
// 采取依赖注入的方式
// 允许用于替换成自定义的Hash函数 也方便测试时替换
// 默认为crc32.ChecksumIEEE算法
type Hash func(data []byte) uint32

// Map是一致性哈希算法的主数据结构
type Map struct {
	// Hash函数
	hash Hash
	// 虚拟节点倍数
	replicas int
	// 哈希环
	keys []int
	// 虚拟节点和真实节点的映射表hashMap
	// key   - 虚拟节点哈希值
	// value - 真实节点名称
	hashMap map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add函数允许传入0或多个真实节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 对于每个真实节点key，对应创建m.replicas个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点的名称为 strconv.Itoa(i)+key
			// 即通过添加编号的方式区分不同虚拟节点
			// 使用m.hash()计算虚拟节点的哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 使用 append 添加到环上
			m.keys = append(m.keys, hash)
			// 在 hashMap 中增加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	// 环上的哈希值排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// 顺时针找到第一个匹配的虚拟节点的下标idx
	// 从m.keys中获取的对应的哈希值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// m.keys是一个环形结构 所以用取余数的方式来处理这种情况
	// 通过hashMap 获得从虚拟节点对应的真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
