package test

import (
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "testing"
)

func TestJsonHelper(t *testing.T) {
    t.Run("object", func(t *testing.T) {
        const JsonObjectText = `{"test_string":"example_string","test_number":114514,"test_bool":true,"test_object":{"some_string":"example_string","some_number":114514,"some_bool":true,"some_array":["some_string",114514,true]},"test_array":["some_string",114514,true,{"some_string":"example_string","some_number":114514,"some_bool":true}]}`

        obj, err := json.DecodeObject(JsonObjectText)
        if err != nil {
            t.Fatal(err)
        }

        if v, err := obj.GetString("test_string"); err != nil || v != "example_string" {
            t.Error(v, err)
        }
        if v, err := obj.GetBool("test_object", "some_array", 2); err != nil || v != true {
            t.Error(v, err)
        }
    })

    t.Run("array", func(t *testing.T) {
        const JsonArrayText = `["some_string",114514,true,{"some_string":"example_string","some_number":114514,"some_bool":true},["some_string",[114514,[true]]]]`
        arr, err := json.DecodeArray(JsonArrayText)
        if err != nil {
            t.Fatal(err)
        }
        if len(arr) != 5 {
            t.Error("len(arr) != 5")
        }
        if v, err := arr.GetBool(4, 1, 1, 0); err != nil || v != true {
            t.Error(v, err)
        }
    })

    t.Run("path_parser", func(t *testing.T) {
        const JsonText = `{"something":{"nothing":114514},"array":[1,1,4,5,1,4]}`
        obj, err := json.DecodeObject(JsonText)
        if err != nil {
            t.Fatal(err)
        }
        if v, err := obj.GetNum(json.ParsePath("/something/nothing")...); err != nil || v != 114514 {
            t.Error(v, err)
        }
        if v, err := obj.GetNum(json.ParsePath("array/2")...); err != nil || v != 4 {
            t.Error(v, err)
        }
    })
}
