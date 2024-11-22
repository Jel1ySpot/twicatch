package api

import (
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "strings"
)

type UserTweets struct {
    PinTweet *Tweet
    Tweets   []Tweet
}

func (u UserTweets) Parse(o *json.JsonObject) UserTweets {
    instructions := o.MustGetArray(json.ParsePath("user/result/timeline_v2/timeline/instructions")...)
    var entries *json.JsonArray = nil
    for i := 0; i < instructions.Length(); i++ {
        switch instructions.MustGetString(i, "type") {
        case "TimelinePinEntry":
            if pinResult, _ := instructions.GetObject(i, "entry", "content", "itemContent", "tweet_results", "result"); pinResult != nil {
                pin := Tweet{}.ParseResult(pinResult)
                u.PinTweet = &pin
            }
        case "TimelineAddEntries":
            entries, _ = instructions.GetArray(i, "entries")
            if entries != nil {
                for i := 0; i < entries.Length(); i++ {
                    t := entries.MustGetObject(i)
                    entryID := t.MustGetString("entryId")
                    switch {
                    case strings.HasPrefix(entryID, "tweet"):
                        u.Tweets = append(u.Tweets, Tweet{}.ParseResult(t.MustGetObject(json.ParsePath("content/itemContent/tweet_results/result")...)))
                    }
                }
            }
        }
    }
    return u
}
