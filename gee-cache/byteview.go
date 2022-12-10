package geecache

type ByteView struct {
	// 储存真实的缓存值
	// byte类型可支持任意数据结构的存储
	b []byte
}

// 实现Len()方法 返回其所占内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// b是只可读的，使用此方法返回一个拷贝，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
