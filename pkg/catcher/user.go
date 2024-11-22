package catcher

import (
    "errors"
    "github.com/Jel1ySpot/twicatch/pkg/api"
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "regexp"
)

const (
    UserByScreenNamePattern = `^https?://(?:twitter|x)\.com/i/api/graphql/[^/]+/UserByScreenName(\?.*)?$`
    UserTweetsPattern       = `^https?://(?:twitter|x)\.com/i/api/graphql/[^/]+/UserTweets(\?.*)?$`
)

func (c *Context) UserTweets(url string) (*api.UserTweets, error) {
    page, err := c.Browser.NewPage()
    if err != nil {
        return nil, err
    }

    var data json.Object

    if page.Context().AddCookies(c.Cookies) != nil {
        return nil, err
    }

    resp, err := page.ExpectResponse(regexp.MustCompile(UserTweetsPattern), func() error {
        _, err := page.Goto(url)
        return err
    })
    if err != nil {
        return nil, err
    }

    if resp.JSON(&data) != nil {
        return nil, err
    }

    _ = page.Close()

    return userTweetsParser((*json.JsonObject)(&data).MustGetObject("data"))
}

func userByScreenNameParser(data *json.JsonObject) (*api.TwitterUser, error) {
    result, err := data.GetObject("result")
    if err != nil {
        return nil, err
    }
    if result.MustGetString("__typename") != "User" {
        return nil, errors.New("data is not user object")
    }
    user := api.TwitterUser{}.Parse(result)
    return &user, nil
}

func userTweetsParser(data *json.JsonObject) (*api.UserTweets, error) {
    instructions := data.MustGetArray(json.ParsePath("user/result/timeline_v2/timeline/instructions")...)
    var entries *json.JsonArray = nil
    for i := 0; i < instructions.Length(); i++ {
        if instructions.MustGetString(i, "type") == "TimelineAddEntries" {
            entries = instructions.MustGetArray(i, "entries")
        }
    }
    if entries == nil {
        return nil, errors.New("parse UserTweets failed")
    }
    userTweets := api.UserTweets{}.ParseEntries(entries)
    return &userTweets, nil
}
