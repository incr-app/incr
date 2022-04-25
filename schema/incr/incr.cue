package incr

import(
    "incr.app/incr/schema/core"
    "strconv"
)

#Version: core.#Version & {
    schemaName: "incr"
    input: int | *1
    output: strconv.FormatInt(input, 10)
}

#Bump: core.#Bump & {
    input: #Version
    output: #Version & { "input": input.input+1 }
}
