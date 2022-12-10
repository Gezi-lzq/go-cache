package singleflight

import "sync"

// 代表进行中/已经结束的请求
type call struct {
	// 使用sync.WaitGroup锁避免重入
	wg  sync.WaitGroup
	val interface{}
	err error
}

// 为singleflight的主数据结构
// 管理不同的Key的请求(call)
type Group struct {
	// 保护Group成员变量m不被并发读写而加上的锁
	mu sync.Mutex
	m  map[string]*call
}

// Do方法
// 针对相同的Key，无论Do被调用多少次，函数fn都只会调用一次
// 等待fn调用结束了，返回返回值或错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	//Do方法会对Group内的m进行读写,因此先上锁
	g.mu.Lock()
	// -------------------------------------------------
	// 饥汉式初始化 用的再初始化
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 该请求已存在，可能正在进行中
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// 请求正在进行中,则等待
		c.wg.Wait()
		// 请求结束，返回结果
		return c.val, c.err
	}
	// 若该请求此时未发过,则创建请求
	c := new(call)
	// 发起请求前上锁
	c.wg.Add(1)
	// 注册进Group.m中,表面Key已经有对应的请求
	g.m[key] = c
	// ------------------------------------------------
	g.mu.Unlock() //不会再对g.m进行读写了,释放锁

	// 调用 fn,发起请求
	c.val, c.err = fn()
	// 请求结束
	c.wg.Done()

	g.mu.Lock() // 下面会对g.m进行修改 因此上锁
	// ------------------------------------------------
	// 更新 g.m
	delete(g.m, key)
	// ------------------------------------------------
	g.mu.Unlock() // 修改完毕，是否锁

	return c.val, c.err // 返回结果
}
