package main

import (
	"math/rand"
	"sync"
)

type UniqueRand struct {
	lock      sync.Mutex
	generated map[uint64]bool
}

func (u *UniqueRand) UInt64() uint64 {
	u.lock.Lock()
	defer u.lock.Unlock()
	for {
		i := rand.Uint64()
		if !u.generated[i] {
			u.generated[i] = true
			return i
		}
	}
}

func (u *UniqueRand) Reset() {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.generated = map[uint64]bool{}
}
