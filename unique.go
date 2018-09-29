package main

import "math/rand"

type UniqueRand struct {
	generated map[uint64]bool
}

func (u *UniqueRand) UInt64() uint64 {
	for {
		i := rand.Uint64()
		if !u.generated[i] {
			u.generated[i] = true
			return i
		}
	}
}

func (u *UniqueRand) Reset() {
	u.generated = map[uint64]bool{}
}
