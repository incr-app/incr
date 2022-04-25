package testcase

import (
	"incr.app/incr/schema/core"
)

#TestComparer: {
	_a:      core.#Version
	_b:      core.#Version
	_output: _a.output & _b.output
}

#TestCases: [testname=string]: #TestCase & {name: testname}

#TestCase: {
	name:   string
	got:    core.#Version
	want:   core.#Version
	verify: #TestComparer & {_a: got, _b: want}
}
