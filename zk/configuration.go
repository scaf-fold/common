package zk

import "time"

// Configuration zookeeper configuration for connecting
type Configuration struct {
	ConnCluster []string // 链接簇
	RootNode    string   // 根节点
	CachePath   string   // 缓存路径
	EndPoint    *Node    // 单前终端信息
}

type Node struct {
	Ip           string // IP 地址
	Port         int64  // 端口
	RegisterTime time.Time
}
