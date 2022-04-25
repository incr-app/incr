package schema

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"incr.app/incr/common"
)

var (
	//go:embed cue.mod/module.cue */*.cue
	embeddedFS embed.FS

	bumpPath       cue.Path = cue.ParsePath("#Bump")
	versionDefPath cue.Path = cue.ParsePath("#Version")
	schemaNamePath cue.Path = cue.ParsePath("schemaName")
	outputPath     cue.Path = cue.ParsePath("output")
	inputPath      cue.Path = cue.ParsePath("input")
	paramsPath     cue.Path = cue.ParsePath("params")
)

func GetEmbeddedSchemas() (map[string]common.Schema, error) {
	result := make(map[string]common.Schema)

	packageNames, err := getCuePackageNames()
	if err != nil {
		return result, err
	}

	ctx := cuecontext.New()

	rawOverlay, err := getOverlay()
	if err != nil {
		return result, err
	}
	for _, packageName := range packageNames {
		td, err := os.MkdirTemp("", "")
		if err != nil {
			return result, err
		}
		defer os.RemoveAll(td)
		overlay := make(map[string]load.Source)
		for rawPath, source := range rawOverlay {
			overlay[path.Join(td, rawPath)] = source
		}
		config := &load.Config{
			Dir:     path.Join(td, packageName),
			Package: packageName,
			Overlay: overlay,
		}
		instance := ctx.BuildInstance(load.Instances(nil, config)[0])
		if err != nil {
			return result, fmt.Errorf("failed to load package %s. Internal error: %w", packageName, err)
		}
		version := instance.LookupPath(versionDefPath)
		bump := instance.LookupPath(bumpPath)
		result[packageName] = &CueSchema{
			version: &version,
			bump:    &bump,
			name:    packageName,
		}
	}

	return result, nil
}

func getOverlay() (map[string]load.Source, error) {
	result := make(map[string]load.Source)
	err := fs.WalkDir(embeddedFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		bytes, err := embeddedFS.ReadFile(p)
		if err != nil {
			return err
		}
		result[p] = load.FromBytes(bytes)
		return nil
	})

	return result, err
}

func getCuePackageNames() ([]string, error) {
	result := make([]string, 0, 2)

	blacklist := []string{
		"core",
		"cue.mod",
		"testcase",
	}

	contains := func(haystack []string, needle string) bool {
		for _, candidate := range haystack {
			if needle == candidate {
				return true
			}
		}
		return false
	}

	dirEntries, err := embeddedFS.ReadDir(".")
	if err != nil {
		return result, err
	}
	for _, dirEntry := range dirEntries {
		name := dirEntry.Name()
		if !contains(blacklist, name) {
			result = append(result, name)
		}
	}

	return result, nil
}

func DebugCueValue(v *cue.Value) {
	nodes := make([]string, 0, 10)
	v.Walk(func(_ cue.Value) bool { return true }, func(v cue.Value) { nodes = append(nodes, v.Path().String()) })
	fmt.Printf("DEBUG: walking %s gave these paths: %s\n", v.Path().String(), strings.Join(nodes, ", "))
	s, err := v.String()
	if err != nil {
		fmt.Printf("DEBUG: but %s could not be stringified because of %s\n", v.Path().String(), err)
	} else {
		fmt.Printf("DEBUG: and %s has value: %s\n", v.Path().String(), s)
	}
}
