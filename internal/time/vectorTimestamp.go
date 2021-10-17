package time

import (
	"strconv"
	"sync"

	math "github.com/goose-alt/chitty-chat/internal/math"
)

type VectorTimestamp struct {
	ClientId   string
	vectorTime map[string]int
	time       int
	lock       sync.Mutex
}

func CreateVectorTimestamp(clientId string) VectorTimestamp {

	return VectorTimestamp{
		ClientId:   clientId,
		vectorTime: make(map[string]int),
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

	for key, vt := range foreignTime.GetVectorTime() {
		
		maxValue := math.Max(v.vectorTime[key], vt)

		v.vectorTime[key] = maxValue
		v.time += maxValue
	}
}

func (v VectorTimestamp) GetVectorTime() map[string]int {
	return v.vectorTime
}

func (v VectorTimestamp) GetDisplayableContent() string {

	v.lock.Lock()
	defer v.lock.Unlock()

	return strconv.Itoa(v.time)
}

func (v VectorTimestamp) Increment() {

	v.lock.Lock()
	defer v.lock.Unlock()

	v.vectorTime[v.ClientId] += 1
	v.time += 1
}
