package json

import (
	"os"
	"sync"

	"incr.app/incr/datastore"
)

type JsonDatastore struct {
	Filename string
}

type jsonFile struct {
	*sync.Mutex
	filename string
}

func NewDatastore(filename string) (JsonDatastore, error) {
	datastore := JsonDatastore{
		Filename: filename,
	}
	if _, err := os.Open(filename); err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(filename, []byte(`{"versions":{}}`), 0600); err != nil {
				return datastore, err
			}
		} else {
			return datastore, err
		}
	}
	return datastore, nil
}

func (jd JsonDatastore) Connect() (datastore.Conn, error) {
	return NewJsonConn(&jsonFile{Mutex: &sync.Mutex{}, filename: jd.Filename})
}
