package semver

import (
	"incr.app/incr/schema/testcase"
)

_stub: #Version & {input: "1.2.3"}
tests: {
	"schema name": #Version.schemaName & "semver"
	cases: testcase.#TestCases & {
		"from structure": {
			got: (#Version & {
				input: {
					major:    1
					minor:    0
					patch:    0
					metadata: "something"
				}
			})
			want: #Version & {input: "1.0.0+something"}
		}

		"from string": {
			got:  #Version & {input: "1.0.1-blah+something"}
			want: #Version & {
				input: {
					major:      1
					minor:      0
					patch:      1
					prerelease: "blah"
					metadata:   "something"
				}
			}
		}
		"bump major": {
			_initial: #Bump & {input: _stub} & {params: type: "major"}
			got:      _initial.output
			want:     #Version & {input: "2.0.0"}
		}
		"bump minor": {
			_initial: #Bump & {input: _stub} & {params: type: "minor"}
			got:      _initial.output
			want:     #Version & {input: "1.3.0"}
		}
		"bump patch": {
			_initial: #Bump & {input: _stub} & {params: type: "patch"}
			got:      _initial.output
			want:     #Version & {input: "1.2.4"}
		}
		"bump with metadata": {
			_initial: #Bump & {input: _stub} & {params: metadata: "metadata"}
			got:      _initial.output
			want:     #Version & {input: "1.2.4+metadata"}
		}
		"bump with prerelease": {
			_initial: #Bump & {input: _stub} & {params: prerelease: "prerelease"}
			got:      _initial.output
			want:     #Version & {input: "1.2.4-prerelease"}
		}
		"bump with both": {
			_initial: #Bump & {input: _stub } & {
				params: {
					prerelease: "prerelease"
					metadata:   "metadata"
				}
			}
			got:  _initial.output
			want: #Version & {input: "1.2.4-prerelease+metadata"}
		}
	}
}
