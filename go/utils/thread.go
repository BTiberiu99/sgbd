package utils

import (
	"sync"
)

//CreateSyncFunc... creates an function with a state of mutex
// that makes any  call reiceived to be in sync
func CreateSyncFunc() func(f func()) {

	mut := new(sync.Mutex)

	return func(f func()) {
		mut.Lock()
		defer mut.Unlock()

		f()
	}
}
