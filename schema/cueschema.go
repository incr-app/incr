package schema

import (
	"fmt"

	"cuelang.org/go/cue"
	"incr.app/incr/common"
)

type CueSchema struct {
	common.Schema
	name    string
	version *cue.Value
	bump    *cue.Value
}

func (s *CueSchema) New(input string) common.Version {
	cv := s.version.FillPath(inputPath, input)
	return NewCueVersion(&cv)
}

func (s *CueSchema) Name() string {
	return s.name
}

func (s *CueSchema) Bump(input common.Version, params common.BumpParams) (common.Version, error) {
	bumper := s.bump.FillPath(cue.ParsePath("input.input"), input.String()).FillPath(paramsPath, params)
	newVersion := bumper.LookupPath(cue.ParsePath("output"))
	// assert that the value is concrete
	if _, err := newVersion.LookupPath(outputPath).String(); err != nil {
		return nil, fmt.Errorf("can't bump version %s. Bumped version is not concrete. Internal error: %w", input, err)
	}
	result := NewCueVersion(&newVersion)

	return result, nil
}

func (s *CueSchema) Debug() {
	DebugCueValue(s.version)
	DebugCueValue(s.bump)
}
