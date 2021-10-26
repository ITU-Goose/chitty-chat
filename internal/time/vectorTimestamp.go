package time

import (
	"strconv"
	"sync"

	math "github.com/goose-alt/chitty-chat/internal/math"
)

type VectorTimestamp struct {
	ClientId   string
	vectorTime map[string]int32
	time       int32
	lock       sync.Mutex
}

func CreateVectorTimestamp(clientId string) VectorTimestamp {

	return VectorTimestamp{
		ClientId:   clientId,
		vectorTime: make(map[string]int32),
		lock:       sync.Mutex{},
	}
}

/*
Synchronizes the two timestamps so that the logical timestamp is updated.
*/
func (v VectorTimestamp) Sync(foreignTime map[string]int32) {

	v.lock.Lock()
	defer v.lock.Unlock()

	v.time = 0 // Reset time and count again

	for key, vt := range foreignTime {

		maxValue := math.Max(v.vectorTime[key], vt)

		v.vectorTime[key] = maxValue
		v.time += maxValue
	}
}

func (v VectorTimestamp) GetVectorTime() map[string]int32 {
	return v.vectorTime
}

func (v VectorTimestamp) GetDisplayableContent() string {

	v.lock.Lock()
	defer v.lock.Unlock()

	return strconv.Itoa(int(v.time))
}

func (v VectorTimestamp) Increment() {

	v.lock.Lock()
	defer v.lock.Unlock()

	v.vectorTime[v.ClientId] += 1
	v.time += 1
}
