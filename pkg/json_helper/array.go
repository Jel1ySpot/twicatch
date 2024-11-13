package json_helper

import (
    "fmt"
)

func (a *JsonArray) Array() *Array {
    return (*Array)(a)
}

func (a *JsonArray) Get(keys ...any) (any, error) {
    if len(keys) == 0 {
        return a, nil
    }

    obj := *a.Array()
    var k int

    if i, ok := keys[0].(int); ok {
        k = i
    } else {
        return nil, fmt.Errorf("key %v is not an int", keys[0])
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

func (a *JsonArray) Length() int {
    return len(*a.Array())
}

func (a *JsonArray) GetObject(keys ...any) (*JsonObject, error) {
    obj, err := a.Get(keys...)
    if err != nil {
        return nil, err
    }
    if v, ok := obj.(Object); ok {
        return (*JsonObject)(&v), nil
    }
    return nil, fmt.Errorf("%v is not an object", obj)
}

func (a *JsonArray) MustGetObject(keys ...any) *JsonObject {
    obj, err := a.GetObject(keys...)
    if err != nil {
        return &JsonObject{}
    }
    return obj
}

func (a *JsonArray) GetArray(keys ...any) (*JsonArray, error) {
    obj, err := a.Get(keys...)
    if err != nil {
        return nil, err
    }
    if v, ok := obj.(Array); ok {
        return (*JsonArray)(&v), nil
    }
    return nil, fmt.Errorf("%v is not an array", obj)
}

func (a *JsonArray) MustGetArray(keys ...any) *JsonArray {
    arr, err := a.GetArray(keys...)
    if err != nil {
        return &JsonArray{}
    }
    return arr
}

func (a *JsonArray) GetNum(keys ...any) (float64, error) {
    obj, err := a.Get(keys...)
    if err != nil {
        return 0, err
    }
    if v, ok := obj.(float64); ok {
        return v, nil
    }
    return 0, fmt.Errorf("%v is not a number", obj)
}

func (a *JsonArray) MustGetNum(keys ...any) float64 {
    num, err := a.GetNum(keys...)
    if err != nil {
        return 0
    }
    return num
}

func (a *JsonArray) GetString(keys ...any) (string, error) {
    obj, err := a.Get(keys...)
    if err != nil {
        return "", err
    }
    if v, ok := obj.(string); ok {
        return v, nil
    }
    return "", fmt.Errorf("%v is not a string", obj)
}

func (a *JsonArray) MustGetString(keys ...any) string {
    str, err := a.GetString(keys...)
    if err != nil {
        return ""
    }
    return str
}

func (a *JsonArray) GetBool(keys ...any) (bool, error) {
    obj, err := a.Get(keys...)
    if err != nil {
        return false, err
    }
    if v, ok := obj.(bool); ok {
        return v, nil
    }
    return false, fmt.Errorf("%v is not a bool", obj)
}

func (a *JsonArray) MustGetBool(keys ...any) bool {
    b, err := a.GetBool(keys...)
    if err != nil {
        return false
    }
    return b
}
