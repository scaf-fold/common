package distribute_id

import (
	"fmt"
	lzk "github.com/scaf-fold/common/zk"
	"log"
	"strconv"
	"sync"
	"time"
)

type DistributeIdImpl struct {
	mutex         sync.Mutex
	workId        int64              // workId
	lastTimestamp *time.Time         // 最后一次时间
	sequence      int64              // 当前序列号
	maxBackOffset int64              // 最大回拨（毫秒）
	zkConf        *lzk.Configuration // zookeeper configuration
}

func NewDistributeIdImpl(start *time.Time, conf *lzk.Configuration) *DistributeIdImpl {
	if start.After(time.Now()) {
		panic("start time must be before")
	}
	return &DistributeIdImpl{
		lastTimestamp: start,
		workId:        0,
		sequence:      0,
		maxBackOffset: 3,
		zkConf:        conf,
	}
}

func (d *DistributeIdImpl) Start() {
	done := make(chan struct{})
	cr, err := lzk.Load(d.zkConf)
	if err != nil {
		panic(err)
	}
	nodeRoot := fmt.Sprintf("/%s/conf", d.zkConf.RootNode)
	log.Println("root node: ", nodeRoot)
	isExists := cr.Exists(nodeRoot)
	if !isExists {
		cr.Start(nodeRoot, done)
	} else {
		cr.Restart(nodeRoot, done)
	}
	<-done
	log.Println("zookeeper has been started")
	d.workId = cr.WorkId()
	if d.workId >= 0 && d.workId <= WorkIdBit.MaxValue() {
		log.Printf("zookeeper initialized successfully,current workId: %d", d.workId)
	} else {
		panic(fmt.Errorf("zookeeper initialized successful,but work id %d is not avaliable,because work id"+
			"must[%d,%d]", d.workId, 0, WorkIdBit.MaxValue()))
	}
	log.Println("Distributed Id Generator started successfully")
}

func (d *DistributeIdImpl) Ids(count int64) ([]string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	curTime := time.Now()
	if d.lastTimestamp != nil && curTime.Before(*d.lastTimestamp) {
		// 时钟回拨
		offset := d.lastTimestamp.UnixMilli() - curTime.UnixMilli()
		if d.maxBackOffset >= offset {
			time.Sleep(time.Millisecond * 6)
			curTime = time.Now()
			if curTime.Before(*d.lastTimestamp) {
				return nil, fmt.Errorf("system Time Call Back")
			}
		} else {
			return nil, fmt.Errorf("system Time Call Back")
		}
	} else if curTime.After(*d.lastTimestamp) {
		d.sequence = 0
	}
	idBuffer := make(chan int64)
	done := make(chan struct{})
	var i int64
	go func(cur time.Time, buffer chan int64, done chan struct{}) {
		for i = 0; i < count; i++ {
			d.sequence = (d.sequence + 1) & SequenceBit.MaxValue()
			if d.sequence == 0 {
				// 下一个周期
				time.Sleep(time.Millisecond)
				cur = time.Now()
			}
			id := (cur.UnixMilli()-d.lastTimestamp.UnixMilli())<<(WorkIdBit+SequenceBit) | d.workId<<SequenceBit | d.sequence
			buffer <- id
		}
		close(idBuffer)
		close(done)
	}(curTime, idBuffer, done)
	result := make([]string, 0)
	for data := range idBuffer {
		result = append(result, strconv.FormatInt(data, 10))
	}
	<-done
	return result, nil
}

func (d *DistributeIdImpl) IdInverse(id string, baseTime time.Time) (*SnowFakeId, error) {
	inverseId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	snowId := &SnowFakeId{}
	// 右移22位位时间戳
	stamp := inverseId>>(WorkIdBit+SequenceBit) + baseTime.UnixMilli()
	snowId.GenerateTime = time.UnixMilli(stamp)
	snowId.WorkId = inverseId>>SequenceBit ^ (inverseId >> (WorkIdBit + SequenceBit) << WorkIdBit)
	snowId.Sequence = inverseId ^ (inverseId >> SequenceBit << SequenceBit)
	snowId.Id = id
	return snowId, nil
}
