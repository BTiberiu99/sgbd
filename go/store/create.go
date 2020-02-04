package store

import (
	"encoding/json"
	"io/ioutil"
)

//GetInstance ... returns the instance of the singleton store
func GetInstance() store {
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

	store := GetInstance()

	_ = json.Unmarshal([]byte(file), &store)

	// fmt.Println(GetInstance())
}
