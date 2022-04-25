package schema

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"incr.app/incr/common"
)

type CueVersion struct {
	common.Version
	version *cue.Value
}

func (cv *CueVersion) String() string {
	serialized, err := cv.version.LookupPath(outputPath).String()
	if err != nil {
		nodes := make([]string, 0, 10)
		cv.version.Walk(func(_ cue.Value) bool { return true }, func(v cue.Value) { nodes = append(nodes, v.Path().String()) })
		fmt.Printf("DEBUG: walking %s gave these paths: %s\n", cv.version.Path().String(), strings.Join(nodes, ", "))
		panic("Can't find serialized version number!")
	}
	return serialized
}

func (cv *CueVersion) SchemaName() string {
	schemaName, err := cv.version.LookupPath(schemaNamePath).String()
	if err != nil {
		panic("Can't find schema name!")
	}
	return schemaName
}

func NewCueVersion(v *cue.Value) *CueVersion {
	return &CueVersion{
		version: v,
	}
}
