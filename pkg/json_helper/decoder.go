package json_helper

import "encoding/json"

func DecodeObject(s string) (JsonObject, error) {
    var v Object
    err := json.Unmarshal([]byte(s), &v)
    return v, err
}

func DecodeArray(s string) (JsonArray, error) {
    var v Array
    err := json.Unmarshal([]byte(s), &v)
    return v, err
}
