package zk

import (
	"errors"
	"github.com/go-zookeeper/zk"
	"time"
)

type Curator struct {
	client      *zk.Conn       // zookeeper 连接器
	workId      int64          // workId
	conf        *Configuration // configuration
	lastUpdated *time.Time     // lastUpdated
}

func NewCurator(client *zk.Conn) *Curator {
	return &Curator{
		client:      client,
		workId:      0,
		conf:        nil,
		lastUpdated: nil,
	}
}

func Load(conf *Configuration) (*Curator, error) {
	if conf != nil && conf.ConnCluster != nil {
		conn, _, err := zk.Connect(conf.ConnCluster, time.Second)
		if err != nil {
			return nil, err
		}
		cr := &Curator{
			client:      conn,
			workId:      0,
			conf:        conf,
			lastUpdated: nil,
		}
		return cr, nil
	}
	return nil, errors.New("not Found Effective Configuration")
}
