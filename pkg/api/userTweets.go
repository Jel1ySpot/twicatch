package api

import (
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "strings"
)

type UserTweets struct {
    Tweets []Tweet
}

func (u UserTweets) ParseEntries(o *json.JsonArray) UserTweets {
    if o == nil {
        return u
    }
    for i := 0; i < o.Length(); i++ {
        t := o.MustGetObject(i)
        entryID := t.MustGetString("entryId")
        switch {
        case strings.HasPrefix(entryID, "tweet"):
            u.Tweets = append(u.Tweets, Tweet{}.ParseResult(t.MustGetObject(json.ParsePath("content/itemContent/tweet_results/result")...)))
        }
    }
    return u
}
