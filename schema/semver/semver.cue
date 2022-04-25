package semver

import (
	"incr.app/incr/schema/core"
	"strconv"
	"regexp"
)

#Bump: core.#Bump & {
	input: #Version
	params: {
		type:        "major" | "minor" | *"patch"
		prerelease?: string
		metadata?:   string
	}
	output: #Version & {
		_metadata: {
			if params.prerelease != _|_ {
				prerelease: params.prerelease
			}
			if params.metadata != _|_ {
				metadata: params.metadata
			}
		}
		if params.type == "major" {
			"input": {
				major: input._version.major + 1
				minor: 0
				patch: 0
			} & _metadata
		}
		if params.type == "minor" {
			"input": {
				major: input._version.major
				minor: input._version.minor + 1
				patch: 0
			} & _metadata
		}
		if params.type == "patch" {
			"input": {
				major: input._version.major
				minor: input._version.minor
				patch: input._version.patch + 1
			} & _metadata
		}
	}
}

#Version: core.#Version & {
	_regex: #"^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<metadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$"#
	schemaName: "semver"
	{
		input: {
			major:       int
			minor:       int
			patch:       int
			metadata?:   string
			prerelease?: string
		}
		_version: input
		output: {
			if _version.metadata == _|_ && _version.prerelease == _|_ {
				"\(_version.major).\(_version.minor).\(_version.patch)"
			}
			if _version.metadata != _|_ && _version.prerelease == _|_ {
				"\(_version.major).\(_version.minor).\(_version.patch)+\(_version.metadata)"
			}
			if _version.metadata == _|_ && _version.prerelease != _|_ {
				"\(_version.major).\(_version.minor).\(_version.patch)-\(_version.prerelease)"
			}
			if _version.metadata != _|_ && _version.prerelease != _|_ {
				"\(_version.major).\(_version.minor).\(_version.patch)-\(_version.prerelease)+\(_version.metadata)"
			}
		}
	} | {
		input:  string
		_match: regexp.FindNamedSubmatch(_regex, input)
		_version: {
			major:      strconv.Atoi(_match."major")
			minor:      strconv.Atoi(_match."minor")
			patch:      strconv.Atoi(_match."patch")
			metadata:   _match."metadata"
			prerelease: _match."prerelease"
		}
		output: input
	}
}
