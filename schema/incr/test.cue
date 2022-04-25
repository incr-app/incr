package incr

import (
	"incr.app/incr/schema/testcase"
)

tests: {
    schemaName: #Version.schemaName & "incr"
    cases: testcase.#TestCases & {
        "create": {
            got:  #Version & {input: 12}
            want: #Version & {input: 12}
        }

        "bump": {
            got:      (#Bump & {input: #Version & {input: 12}}).output
            want:     #Version & {input: 13}
        }
    }
}
