package store

import (
	"encoding/json"
	"io/ioutil"
	"sgbd4/go/db"
	"sgbd4/go/utils"
)

const (
	fileNameStore = "store.json"
)

var (
	instance store
)

type store map[string]db.Connection
type safeStore map[string][]byte

func (s store) MarshalJSON() ([]byte, error) {
	safeStore := make(safeStore)

	var (
		err error
	)
	runThreadSafe(func() {
		var m []byte
		for i, val := range s {
			m, err = json.Marshal(val)
			if err != nil {
				return
			}
			safeStore[i] = utils.Encrypt(m)
		}
	})

	if err != nil {
		return []byte{}, nil
	}

	return json.Marshal(safeStore)
}

func (s store) UnmarshalJSON(data []byte) error {
	safeStore := make(safeStore)

	err := json.Unmarshal(data, &safeStore)

	if err != nil {
		return err
	}
	runThreadSafe(func() {
		for i, val := range safeStore {
			conn := new(db.Connection)
			json.Unmarshal(utils.Decrypt(val), &conn)
			s[i] = *conn
		}
	})

	return nil
}
func (s store) Add(conn db.Connection) {
	runThreadSafe(func() {
		s[conn.SafeString()] = conn
	})
	s.save()
}

func (s store) Get(index string) (db.Connection, bool) {
	var (
		conn  db.Connection
		exist bool
	)

	runThreadSafe(func() {
		conn, exist = s[index]
	})

	return conn, exist
}

func (s store) Remove(index string) {
	runThreadSafe(func() {
		delete(s, index)
	})
	s.save()
}

func (s store) Connections() map[string]db.Connection {
	z := make(map[string]db.Connection)

	runThreadSafe(func() {
		for i := range s {
			z[i] = s[i]
		}
	})

	return z
}

func (s store) save() {
	file, _ := json.MarshalIndent(s, "", " ")

	_ = ioutil.WriteFile(fileNameStore, file, 0644)
}
