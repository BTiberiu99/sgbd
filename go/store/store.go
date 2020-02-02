package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sgbd4/go/db"
	"sgbd4/go/utils"
)

const (
	fileNameStore = "store.json"
	activeIndex   = "ai"
)

var (
	instance      store
	runThreadSafe = utils.CreateSyncFunc()
)

type store map[string]db.Connection
type safeStore map[string][]byte

func (s store) MarshalJSON() ([]byte, error) {
	safeStore := make(safeStore)

	for i, val := range s {
		m, err := json.Marshal(val)
		if err != nil {
			return []byte{}, nil
		}

		safeStore[i] = utils.Encrypt(m)
	}

	safeStore[activeIndex] = []byte(db.ActiveIndex)

	return json.Marshal(safeStore)
}

func (s store) UnmarshalJSON(data []byte) error {
	safeStore := make(safeStore)

	err := json.Unmarshal(data, &safeStore)

	if err != nil {
		return err
	}

	for i, val := range safeStore {
		if i == activeIndex {
			db.ActiveIndex = string(val)
			continue
		}
		conn := new(db.Connection)
		json.Unmarshal(utils.Decrypt(val), &conn)
		s[i] = *conn
	}

	//Init connection
	if copy, exist := s[db.ActiveIndex]; exist {
		db.UpdateConnection(&copy)
	}

	return nil
}

//Add ... add new connection to the store
func (s store) Add(conn db.Connection) {
	runThreadSafe(func() {
		s[conn.SafeString()] = conn
	})
	s.writeToDisk()
}

//Get ... get the connection from the store
func (s store) Get(index string) (db.Connection, bool) {
	var (
		conn  db.Connection
		exist bool
	)

	runThreadSafe(func() {
		conn, exist = s[index]
		fmt.Println(s)
	})

	return conn, exist
}

//Remove ... remove the connection from the store
func (s store) Remove(index string) {

	runThreadSafe(func() {
		delete(s, index)
	})

	s.writeToDisk()
}

//Connections ... takes all connections from the store
func (s store) Connections() map[string]db.Connection {
	z := make(map[string]db.Connection)

	runThreadSafe(func() {
		for i := range s {
			z[i] = s[i]
		}
	})

	return z
}

//Save ... save the store to the hard disk
func (s store) Save() {
	s.writeToDisk()
}

func (s *store) writeToDisk() {
	runThreadSafe(func() {
		file, _ := json.MarshalIndent(s, "", " ")

		_ = ioutil.WriteFile(fileNameStore, file, 0644)
	})
}
