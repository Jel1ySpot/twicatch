package json_helper

import (
    "strconv"
    "strings"
)

func ParsePath(path string) []any {
    var result []any

    for _, k := range strings.Split(path, "/") {
        if n, err := strconv.ParseInt(k, 10, 32); err == nil {
            result = append(result, int(n))
        } else {
            result = append(result, k)
        }
    }

    if result[0] == "" {
        result = result[1:]
    }

    return result
}
