package catcher

import (
    "errors"
    "github.com/Jel1ySpot/twicatch/pkg/api"
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "regexp"
)

const (
    TweetDetailPattern = `^https?://(?:twitter|x)\.com/i/api/graphql/[^/]+/TweetDetail(\?.*)?$`
    TweetPattern       = `^https?://api\.(?:twitter|x)\.com/graphql/[^/]+/TweetResultByRestId(\?.*)?$`
)

func (c *Context) Status(url string) (*api.Tweet, error) {
    page, err := c.Browser.NewPage()
    if err != nil {
        return nil, err
    }

    var (
        tweetMatch = regexp.MustCompile(TweetDetailPattern + "|" + TweetPattern)
        data       json.Object
    )

    if page.Context().AddCookies(c.Cookies) != nil {
        return nil, err
    }

    resp, err := page.ExpectResponse(tweetMatch, func() error {
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

    return tweetParser((*json.JsonObject)(&data).MustGetObject("data"))
}

func tweetParser(data *json.JsonObject) (*api.Tweet, error) {
    if _, err := data.GetObject("tweetResult"); err == nil {
        result, err := data.GetObject("tweetResult", "result")
        if err != nil {
            return nil, err
        }
        if result.MustGetString("__typename") == "TweetUnavailable" {
            s, err := result.GetString("reason")
            if err != nil {
                return nil, err
            }
            return nil, errors.New(s)
        }
        tweet := api.Tweet{}.ParseResult(result)
        return &tweet, nil
    } else if instructions, err := data.GetArray("threaded_conversation_with_injections_v2", "instructions"); err == nil {
        tweet := api.Tweet{}.ParseInstructions(instructions)
        return &tweet, nil
    } else {
        return nil, errors.New("parse tweet failed")
    }
}
