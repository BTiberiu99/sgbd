package utils

import (
	"sync"
)

func CreateSyncFunc() func(f func()) {

	mut := new(sync.Mutex)

	return func(f func()) {

		// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		// defer cancel()

		mut.Lock()

		// fmt.Println("Am ajuns aici")
		defer mut.Unlock()
		f()
		// ctx.Done()

		// fmt.Println("Am plecat de aici")
	}
}
