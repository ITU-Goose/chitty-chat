package time

import (
	"strconv"
	"sync"

	math "github.com/goose-alt/chitty-chat/internal/math"
)

type VectorTimestamp struct {
	ClientId   string
	VectorTime map[string]int
	time       int
	lock       sync.Mutex
}

func CreateVectorTimestamp(clientId string) VectorTimestamp {

	return VectorTimestamp{
		ClientId:   clientId,
		VectorTime: make(map[string]int),
		lock:       sync.Mutex{},
	}
}

/*
Synchronizes the two timestamps so that the logical timestamp is updated.
*/
func (v VectorTimestamp) Sync(foreignTime VectorTimestamp) {

	v.lock.Lock()
	defer v.lock.Unlock()

	v.time = 0 // Reset time and count again

	for key, vt := range foreignTime.VectorTime {
		
		maxValue := math.Max(v.VectorTime[key], vt)

		v.VectorTime[key] = maxValue
		v.time += maxValue
	}
}

/*
Returns the timestamp as a displayble string.
*/
func (v VectorTimestamp) GetDisplayableContent() string {

	v.lock.Lock()
	defer v.lock.Unlock()

	return strconv.Itoa(v.time)
}

/*
Increments the timestamp by 1
*/
func (v VectorTimestamp) Increment() {

	v.lock.Lock()
	defer v.lock.Unlock()

	v.VectorTime[v.ClientId] += 1
	v.time += 1
}
