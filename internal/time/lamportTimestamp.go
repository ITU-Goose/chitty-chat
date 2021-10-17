package time

import (
	"strconv"
	"sync"

	math "github.com/goose-alt/chitty-chat/internal/math"
)

var (
	instance *LamportTimestamp
)

type LamportTimestamp struct {
	time int
	lock sync.Mutex
}

func GetLamportTimeStamp() *LamportTimestamp {

	if instance == nil {
		timestamp := LamportTimestamp{time: 0, lock: sync.Mutex{}}
		instance = &timestamp
	}

	return instance
}

func (v *LamportTimestamp) Sync(timestamp LamportTimestamp) {
	v.lock.Lock()
	v.time = math.Max(timestamp.time, v.time)
	v.lock.Unlock()
}

func (v *LamportTimestamp) GetDisplayableContent() string {
	return strconv.Itoa(v.time)
}

func (v *LamportTimestamp) Increment() {
	v.lock.Lock()
	v.time += 1
	v.lock.Unlock()
}
