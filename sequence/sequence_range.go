package sequence

import "sync/atomic"

type SequenceRange struct {
	max   int64
	min   int64
	value int64
	over  atomic.Value
}

func NewSequenceRange(min, max int64) *SequenceRange {
	return &SequenceRange{min: min, max: max, value: min}
}

func (sr *SequenceRange) GetAndIncrement() int64 {
	currentValue := atomic.AddInt64(&sr.value, 1)
	if currentValue > sr.max {
		sr.over.Store(true)
		return -1
	}
	return currentValue
}

func (sr *SequenceRange) IsOver() bool {
	x := sr.over.Load()
	over, _ := x.(bool)
	return over
}
