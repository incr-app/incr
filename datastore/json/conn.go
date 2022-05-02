package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"incr.app/incr/common"
	"incr.app/incr/datastore"
)

type jsonConn struct {
	*jsonFile
	cache *jsonContents
}

func (conn *jsonConn) Get(id datastore.Id) (common.Version, error) {
	v, ok := conn.cache.Versions[id]
	if !ok {
		conn.populateCache()
		v, ok = conn.cache.Versions[id]
		if !ok {
			return nil, datastore.VersionNotFoundError(fmt.Errorf("could not find version with id %s", id))
		}
	}
	return v, nil
}

func (conn *jsonConn) Put(id datastore.Id, version common.Version) error {
	jsonVersion := &JsonVersion{
		Version: version.String(),
		Schema:  version.SchemaName(),
	}
	conn.Lock()
	defer conn.Unlock()
	conn.populateCache()
	conn.cache.Versions[id] = jsonVersion
	contentBytes, err := json.Marshal(conn.cache)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(conn.filename, contentBytes, 0600)
	return err
}

// Close has no real purpose for a file store
func (conn *jsonConn) Close() error { return nil }

func (conn *jsonConn) populateCache() error {
	var contents jsonContents

	contentBytes, err := ioutil.ReadFile(conn.filename)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(contentBytes, &contents); err != nil {
		return err
	}
	for id, versioninfo := range contents.Versions {
		conn.cache.Versions[id] = &JsonVersion{
			Schema:  versioninfo.Schema,
			Version: versioninfo.Version,
		}
	}
	return nil
}

func NewJsonConn(jsonFile *jsonFile) (*jsonConn, error) {
	cache := &jsonContents{
		Versions: make(map[datastore.Id]*JsonVersion),
	}
	jc := &jsonConn{jsonFile: jsonFile, cache: cache}
	jc.Lock()
	defer jc.Unlock()
	err := jc.populateCache()

	return jc, err
}

type jsonContents struct {
	Versions map[datastore.Id]*JsonVersion
}
