package json_helper

import (
    "fmt"
)

func (o *JsonObject) Object() *Object {
    return (*Object)(o)
}

func (o *JsonObject) Keys() []string {
    obj := o.Object()
    var result []string
    for s, _ := range *obj {
        result = append(result, s)
    }
    return result
}

func (o *JsonObject) Get(keys ...any) (any, error) {
    if len(keys) == 0 {
        return o, nil
    }

    obj := *o.Object()
    var k string

    if s, ok := keys[0].(string); ok {
        k = s
    } else {
        return nil, fmt.Errorf("key %v is not a string", keys[0])
    }

    if len(keys) == 1 {
        return obj[k], nil
    }

    switch v := obj[k].(type) {
    case Object:
        return (*JsonObject)(&v).Get(keys[1:]...)
    case Array:
        return (*JsonArray)(&v).Get(keys[1:]...)
    default:
        return nil, fmt.Errorf("value in %v is unable to index ", k)
    }
}

func (o *JsonObject) GetObject(keys ...any) (*JsonObject, error) {
    obj, err := o.Get(keys...)
    if err != nil {
        return nil, err
    }
    if v, ok := obj.(Object); ok {
        return (*JsonObject)(&v), nil
    }
    return nil, fmt.Errorf("%v is not an object", obj)
}

func (o *JsonObject) MustGetObject(keys ...any) *JsonObject {
    obj, err := o.GetObject(keys...)
    if err != nil {
        return &JsonObject{}
    }
    return obj
}

func (o *JsonObject) GetArray(keys ...any) (*JsonArray, error) {
    obj, err := o.Get(keys...)
    if err != nil {
        return nil, err
    }
    if v, ok := obj.(Array); ok {
        return (*JsonArray)(&v), nil
    }
    return nil, fmt.Errorf("%v is not an array", obj)
}

func (o *JsonObject) MustGetArray(keys ...any) *JsonArray {
    arr, err := o.GetArray(keys...)
    if err != nil {
        return &JsonArray{}
    }
    return arr
}

func (o *JsonObject) GetNum(keys ...any) (float64, error) {
    obj, err := o.Get(keys...)
    if err != nil {
        return 0, err
    }
    if v, ok := obj.(float64); ok {
        return v, nil
    }
    return 0, fmt.Errorf("%v is not a number", obj)
}

func (o *JsonObject) MustGetNum(keys ...any) float64 {
    num, err := o.GetNum(keys...)
    if err != nil {
        return 0
    }
    return num
}

func (o *JsonObject) GetString(keys ...any) (string, error) {
    obj, err := o.Get(keys...)
    if err != nil {
        return "", err
    }
    if v, ok := obj.(string); ok {
        return v, nil
    }
    return "", fmt.Errorf("%v is not a string", obj)
}

func (o *JsonObject) MustGetString(keys ...any) string {
    str, err := o.GetString(keys...)
    if err != nil {
        return ""
    }
    return str
}

func (o *JsonObject) GetBool(keys ...any) (bool, error) {
    obj, err := o.Get(keys...)
    if err != nil {
        return false, err
    }
    if v, ok := obj.(bool); ok {
        return v, nil
    }
    return false, fmt.Errorf("%v is not a bool", obj)
}

func (o *JsonObject) MustGetBool(keys ...any) bool {
    b, err := o.GetBool(keys...)
    if err != nil {
        return false
    }
    return b
}
