package catcher

import (
    "errors"
    "github.com/Jel1ySpot/twicatch/pkg/api"
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "github.com/playwright-community/playwright-go"
    "regexp"
    "time"
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

    tweetDetailMatch := regexp.MustCompile(TweetDetailPattern)
    tweetMatch := regexp.MustCompile(TweetPattern)

    var (
        done    = make(chan any)
        data    json.Object
        timeout bool
    )

    page.OnResponse(func(rp playwright.Response) {
        go func() {
            if tweetMatch.MatchString(rp.URL()) || tweetDetailMatch.MatchString(rp.URL()) {
                err = rp.JSON(&data)
                close(done)
            }
        }()
    })

    ctx := page.Context()
    if err = ctx.AddCookies(c.Cookies); err != nil {
        return nil, err
    }

    if _, err = page.Goto(url); err != nil {
        return nil, err
    }

    t := time.AfterFunc(30*time.Second, func() {
        timeout = true
        close(done)
    })

    <-done
    if !timeout {
        t.Stop()
    }
    page.Close()
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
