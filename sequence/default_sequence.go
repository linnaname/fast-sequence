package sequence

import (
	"errors"
	"sync"
)

type DefaultSequence struct {
	sync.Mutex
	sequenceRangeMap map[string]*SequenceRange
}

func (dseq *DefaultSequence) NextId(name string) (nextId int64, err error) {
	dseq.Mutex.Lock()
	currentRange, ok := dseq.sequenceRangeMap[name]
	if !ok {
		newRange, err := NextRange(name)
		if err != nil {
			return -1, nil
		}
		dseq.sequenceRangeMap[name] = newRange
		currentRange = newRange
	}
	dseq.Mutex.Unlock()

	value := currentRange.GetAndIncrement()
	if value == -1 {
		dseq.Mutex.Lock()
		for {
			if currentRange.IsOver() {
				currentRange, err := NextRange(name)
				if err != nil {
					break
				}
				value = currentRange.GetAndIncrement()
				if value != -1 {
					break
				}
			}
		}
		dseq.Mutex.Unlock()
	}

	if value < 0 {
		return -1, errors.New("Sequence value overflow")
	}
	return value, nil
}
