package store

import (
	"encoding/json"
	"io/ioutil"
)

//NewStore ... returns the instance of the singleton store
func NewStore() store {
	runThreadSafe(func() {
		if instance == nil {
			instance = make(store)
		}
	})
	return instance
}

func init() {

	file, err := ioutil.ReadFile(fileNameStore)
	if err != nil {
		return
	}

	store := NewStore()

	_ = json.Unmarshal([]byte(file), &store)

	// fmt.Println(NewStore())
}
