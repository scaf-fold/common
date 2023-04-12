package distribute_id

import "time"

type BitCount int64

const (
	SymbolBit      BitCount = 1  // symbol location bit count
	MillisecondBit BitCount = 41 // milliseconds bit count
	WorkIdBit      BitCount = 10 // work id bits count
	SequenceBit    BitCount = 12 // sequence bits count
)

func (b BitCount) MaxValue() int64 {
	return ^(-1 << b)
}

type SnowFakeId struct {
	Id           string
	GenerateTime time.Time
	WorkId       int64
	Sequence     int64
}

/*
	|符号位|时间戳｜工作区ID|序列号|

1）总长度位64位；
2）从第二位开始使用，第一位为符号位，其余位为数据位；
*/
// DistributedId 分布式ID
type DistributedId interface {
	Ids(count int64) ([]string, error)
	IdInverse(id string, baseTime time.Time) (*SnowFakeId, error)
}
