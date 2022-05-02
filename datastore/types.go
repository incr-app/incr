package datastore

import (
	"incr.app/incr/common"
)

type Id string

type Datastore interface {
	// Open a connection to the datastore
	Connect() (Conn, error)
}

type Conn interface {
	// Get a version by ID
	Get(id Id) (common.Version, error)
	// Write a version by ID
	Put(id Id, version common.Version) error
	// Close the connection
	Close() error
}

type VersionNotFoundError error

func IsVersionNotFoundError(err error) bool {
	switch err.(type) {
	case VersionNotFoundError:
		return true
	default:
		return false
	}
}
