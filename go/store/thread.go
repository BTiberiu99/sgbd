package store

import "sync"

var (
	lock = sync.Mutex{}
)

func runThreadSafe(f func()) {
	lock.Lock()
	defer lock.Unlock()
	f()
}
