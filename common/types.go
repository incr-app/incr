package common

import "fmt"

type BumpParams map[string]interface{}

type Schema interface {
	Name() string
	Bump(Version, BumpParams) (Version, error)
	New(input string) Version
}

type Version interface {
	fmt.Stringer
	SchemaName() string
}
