package zk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (cr *Curator) CreateParentsIfNeededWithCustomPolicy(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	var nodePath string
	if path == "" {
		return string(filepath.Separator), nil
	}
	paths := strings.Split(path, string(filepath.Separator))
	size := len(paths)
	for index, sub := range paths {
		log.Println("size", size, ", index", index, ", sub", sub)
		if sub == "" {
			continue
		}
		nodePath += fmt.Sprintf("%c/%s", filepath.Separator, sub)
		flag, _, _ := cr.client.Exists(nodePath)
		if flag {
			continue
		}
		if index < size-1 {
			nodePath, _ = cr.client.Create(nodePath, []byte(""), 0, acl)
			log.Println("index=", index, "nodePath=", nodePath)
		} else {
			nodePath, _ = cr.client.Create(nodePath, data, flags, acl)
			log.Println("index=", index, "nodePath=", nodePath)
		}
	}
	return nodePath, nil
}

func (cr *Curator) DefaultNode() ([]byte, error) {
	if cr.conf.EndPoint != nil {
		data := cr.conf.EndPoint
		data.RegisterTime = time.Now()
		return json.Marshal(data)
	}
	return nil, errors.New("not found for endpoint")
}

// Add 添加node key=path/ip:port-workId
func (cr *Curator) Add(nodePath string, buildData func() ([]byte, error)) (string, error) {
	nodeKey := fmt.Sprintf("%s/%s:%d-", nodePath, cr.conf.EndPoint.Ip, cr.conf.EndPoint.Port)
	var data []byte
	var err error
	if buildData == nil {
		data, err = cr.DefaultNode()
	} else {
		data, err = buildData()
	}
	if err != nil {
		return "", err
	}
	return cr.CreateParentsIfNeededWithCustomPolicy(nodeKey, data, zk.FlagSequence, zk.WorldACL(zk.PermAll))
}

func (cr *Curator) MountSequenceId() {
	_, err := os.Stat(cr.conf.CachePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(cr.conf.CachePath), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	c, err := json.Marshal(cr.workId)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(cr.conf.CachePath, c, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (cr *Curator) Set(key string, data []byte) (*zk.Stat, error) {
	flag, stat, err := cr.client.Exists(key)
	if err == nil && flag {
		return cr.client.Set(key, data, stat.Version)
	}
	return stat, err
}

// Check 当前节点是否时钟回拨
func (cr *Curator) Check(nodeKey string) {
	data, _, err := cr.client.Get(nodeKey)
	if err != nil {
		panic(err)
	}

	lastNodeInfo := &Node{}
	err = json.Unmarshal(data, lastNodeInfo)
	if err != nil {
		panic(err)
	}

	if lastNodeInfo.RegisterTime.After(time.Now()) {
		panic(fmt.Errorf("current time is invalid, last time is %v", lastNodeInfo.RegisterTime))
	}
}

func (cr *Curator) ScheduleNode(node string) {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t := time.Now()
			cr.lastUpdated = &t
			data, err := cr.DefaultNode()
			if err != nil {
				fmt.Println(err)
				break
			}
			_, err = cr.Set(node, data)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func (cr *Curator) Exists(node string) bool {
	flag, _, _ := cr.client.Exists(node)
	return flag
}

func (cr *Curator) Children(node string) ([]string, *zk.Stat, error) {
	return cr.client.Children(node)
}

func (cr *Curator) WorkId() int64 {
	return cr.workId
}

func (cr *Curator) Start(root string, done chan struct{}) {
	defer close(done)
	node, err := cr.Add(root, nil)
	if err != nil {
		panic(err)
	}
	cr.MountSequenceId()
	go cr.ScheduleNode(node)
}

func (cr *Curator) Restart(root string, done chan struct{}) {
	defer close(done)
	nodeMap := make(map[string]int64)
	realMap := make(map[string]string)
	keys, _, err := cr.Children(root)
	if err != nil {
		fmt.Println(err)
	}
	if len(keys) > 0 {
		fmt.Println(keys)
		for _, key := range keys {
			nodeKey := strings.Split(key, "-")
			if len(nodeKey) > 0 {
				realMap[nodeKey[0]] = key
				if len(nodeKey) > 1 {
					i, err := strconv.ParseInt(nodeKey[1], 10, 64)
					if err != nil {
						i = -1
					}
					nodeMap[nodeKey[0]] = i
				}
			}
			listenAddr := fmt.Sprintf("%s:%d", cr.conf.EndPoint.Ip, cr.conf.EndPoint.Port)
			var workId int64 = -1
			if realId, ok := nodeMap[listenAddr]; ok {
				workId = realId
			}
			fmt.Println("workId = ", workId)
			nodeAddr := root + "/" + realMap[listenAddr]
			if workId != -1 {
				cr.workId = workId
				cr.Check(nodeAddr)
				go cr.ScheduleNode(nodeAddr)
				cr.MountSequenceId()
			} else {
				p, err := cr.Add(root, nil)
				if err != nil {
					panic(err)
				}
				sub := strings.ReplaceAll(p, filepath.Dir(p), "")
				nodeKey := strings.Split(sub, "-")
				i, err := strconv.ParseInt(nodeKey[1], 10, 64)
				if err != nil {
					panic(err)
				}
				cr.workId = i
				cr.MountSequenceId()
			}
		}
	} else {
		p, err := cr.Add(root, nil)
		if err != nil {
			panic(err)
		}
		sub := strings.ReplaceAll(p, filepath.Dir(p), "")
		nodeKey := strings.Split(sub, "-")
		i, err := strconv.ParseInt(nodeKey[1], 10, 64)
		if err != nil {
			panic(err)
		}
		cr.workId = i
		cr.MountSequenceId()
	}

}
