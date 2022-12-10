package main

import (
    "flag"
    "fmt"
    "geecache"
    "log"
    "net/http"
)

// Overall flow char										     requsets					        local
// gee := createGroup() --------> /api Service : 9999 ------------------------------> gee.Get(key) ------> g.mainCache.Get(key)
// 						|											^					|
// 						|											|					|remote
// 						v											|					v
// 				cache Service : 800x								|			g.peers.PickPeer(key)
// 						|create hash ring & init peerGetter			|					|
// 						|registry peers write in g.peer				|					|p.httpGetters[p.hashRing(key)]
// 						v											|					|
//			httpPool.Set(otherAddrs...)								|					v
// 		g.peers = gee.RegisterPeers(httpPool)						|			g.getFromPeer(peerGetter, key)
// 						|											|					|
// 						|											|					|
// 						v											|					v
// 		http.ListenAndServe("localhost:800x", httpPool)<------------+--------------peerGetter.Get(key)
// 						|											|
// 						|requsets									|
// 						v											|
// 					p.ServeHttp(w, r)								|
// 						|											|
// 						|url.parse()								|
// 						|--------------------------------------------

// 使用 map 模拟了数据源 db
var db = map[string]string{
    "Tom":  "630",
    "Jack": "589",
    "Sam":  "567",
}

func createGroup() *geecache.Group {
    return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
        func(key string) ([]byte, error) {
            log.Println("[SlowDB] search key", key)
            if v, ok := db[key]; ok {
                return []byte(v), nil
            }
            return nil, fmt.Errorf("%s not exist", key)
        }))
    }

    // 用来启动缓存服务器
    func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
        // 创建HTTPPool
        peers := geecache.NewHTTPPool(addr)
        // 添加节点信息
        peers.Set(addrs...)
        // 注册到gee中
        gee.RegisterPeers(peers)
        log.Println("geecache is running at", addr)
        log.Fatal(http.ListenAndServe(addr[7:], peers))
    }

    // 启动一个API服务 与用户进行交互 用户感知
    func startAPIServer(apiAddr string, gee *geecache.Group) {
        http.Handle("/api", http.HandlerFunc(
            func(w http.ResponseWriter, r *http.Request) {
                key := r.URL.Query().Get("key")
                view, err := gee.Get(key)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                w.Header().Set("Content-Type", "application/octer-stream")
                w.Write(view.ByteSlice())
            }))
            log.Println("fontend server is running at", apiAddr)
            log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
        }

        func main() {
            var port int
            var api bool
            // main()函数需要命令行传入port和api 2个参数
            // 用来在指定端口启动HTTP服务
            flag.IntVar(&port, "port", 8001, "Geecache server port")
            flag.BoolVar(&api, "api", false, "Start a api server?")
            flag.Parse()

            apiAddr := "http://localhost:9999"
            addrMap := map[int]string{
                8001: "http://localhost:8001",
                8002: "http://localhost:8002",
                8003: "http://localhost:8003",
            }

            var addrs []string
            for _, v := range addrMap {
                addrs = append(addrs, v)
            }
            gee := createGroup()
            if api {
                // startAPIServer() 用来启动一个 API 服务（端口 9999）
                // 与用户进行交互，用户感知
                go startAPIServer(apiAddr, gee)
            }
            // 用来启动缓存服务器：创建 HTTPPool，添加节点信息，注册到 gee 中
            // 启动 HTTP 服务（共3个端口，8001/8002/8003），用户不感知
            startCacheServer(addrMap[port], []string(addrs), gee)
        }
