package sequence

import (
	"github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
)

const SIZE = 100

var gobalSQ = &DefaultSequence{sequenceRangeMap: make(map[string]*SequenceRange, 5)}

func TestNextId(t *testing.T) {
	sq := &DefaultSequence{sequenceRangeMap: make(map[string]*SequenceRange, 5)}
	id, err := sq.NextId("test")
	assert.NoError(t, err)
	assert.Greater(t, id, int64(1))
	nextId, err := sq.NextId("test")
	assert.Greater(t, nextId, id)
	assert.Equal(t, nextId, id+1)
}

func TestSerialGenerate(t *testing.T) {
	sq := &DefaultSequence{sequenceRangeMap: make(map[string]*SequenceRange, 5)}
	uidSet := mapset.NewThreadUnsafeSet()
	for i := 0; i < SIZE; i++ {
		doGenerate(t, sq, uidSet, i)
	}
	checkUniqueID(t, uidSet)
}

func TestParallelGenerate(t *testing.T) {
	uidSet := mapset.NewSet()
	sq := &DefaultSequence{sequenceRangeMap: make(map[string]*SequenceRange, 5)}
	var wg sync.WaitGroup
	control := int32(-1)
	var mutex sync.Mutex

	for i := 0; i < runtime.NumCPU()<<2; i++ {
		wg.Add(1)
		go func() {
			mutex.Lock()
			for {
				if control != SIZE {
					control = control + 1
				}
				if control == SIZE {
					break
				}
				doGenerate(t, sq, uidSet, int(control))
			}
			mutex.Unlock()
			wg.Done()
		}()

	}
	wg.Wait()
	assert.Equal(t, control, int32(SIZE))
	checkUniqueID(t, uidSet)
}

func doGenerate(t *testing.T, sq *DefaultSequence, uidSet mapset.Set, index int) {
	id, err := sq.NextId("test")
	assert.NoError(t, err)
	assert.True(t, id > int64(0))
	uidSet.Add(id)
}

func checkUniqueID(t *testing.T, uidSet mapset.Set) {
	assert.Equal(t, SIZE, uidSet.Cardinality())
}

func BenchmarkNextId(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gobalSQ.NextId("test")
	}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var j = 0
		j++
	}
}
