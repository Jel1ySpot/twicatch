package json_helper

import "encoding/json"

func (o *JsonObject) Encode(indent string) (string, error) {
    if indent == "" {
        b, err := json.Marshal(o)
        return string(b), err
    }
    b, err := json.MarshalIndent(o, "", indent)
    return string(b), err
}

func (a *JsonArray) Encode(indent string) (string, error) {
    if indent == "" {
        b, err := json.Marshal(a)
        return string(b), err
    }
    b, err := json.MarshalIndent(a, "", indent)
    return string(b), err
}
