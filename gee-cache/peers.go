package geecache

import pd "geecache/geecachepb"

// PickPeer()方法用于根据传入Key选择相应的节点PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// Get()方法用于从对应的group查找缓存值
// PeerGetter即为对应的流程中的HTTP客户端
type PeerGetter interface {
	Get(in *pd.Request, out *pd.Response) error
}
