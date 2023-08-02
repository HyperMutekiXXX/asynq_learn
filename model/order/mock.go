package order

import "sync/atomic"

var AtId int32

var OrderMap = make(map[int32]*Order)

func GetAtId() int32 {
	atomic.AddInt32(&AtId, 1)
	return AtId
}
