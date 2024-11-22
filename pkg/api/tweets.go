package api

import (
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "strconv"
    "strings"
    "time"
)

const (
    Photo TwitterEntityMediaType = iota
    Video
    Gif

    TimeForm = `Mon Jan 02 15:04:05 -0700 2006`
)

type (
    TwitterEntityMediaType int

    TwitterEntityMedia struct {
        Type        TwitterEntityMediaType
        MediaUrl    string
        VideoUrl    []string
        ExpandedUrl string
    }

    TwitterEntityUrl struct {
        DisplayUrl  string
        ExpandedUrl string
        Url         string
    }

    TwitterEntityUserMention struct {
        ID         string
        Name       string
        ScreenName string
    }

    TwitterEntities struct {
        HashTag      []string
        Media        []TwitterEntityMedia
        Symbols      []string
        Timestamps   []time.Time //TODO
        Urls         []TwitterEntityUrl
        UserMentions []TwitterEntityUserMention
    }

    Tweet struct {
        RestID              string
        CreatedAt           time.Time
        FullText            string
        Lang                string
        PossiblySensitive   bool
        Entities            TwitterEntities
        ConversationThreads [][]Tweet //TODO
        User                TwitterUser
        ViewsCount          int64
        BookmarkCount       int
        FavoriteCount       int
        QuoteCount          int
        ReplyCount          int
        RetweetCount        int
        Bookmarked          bool
        Favorited           bool
        Retweeted           bool
    }
)

func (t TwitterEntityUserMention) Parse(o *json.JsonObject) TwitterEntityUserMention {
    if o == nil {
        return t
    }
    t.ID = o.MustGetString("id_str")
    t.Name = o.MustGetString("name")
    t.ScreenName = o.MustGetString("screen_name")
    return t
}

func (t TwitterEntityMedia) Parse(o *json.JsonObject) TwitterEntityMedia {
    if o == nil {
        return t
    }
    s := o.MustGetString("type")
    switch s {
    case "photo":
        t.Type = Photo
    case "video":
        t.Type = Video
    case "animated_gif":
        t.Type = Gif
    }
    t.MediaUrl = o.MustGetString("media_url_https")
    if v, err := o.GetArray("video_info", "variants"); err == nil {
        for i := 0; i < v.Length(); i++ {
            t.VideoUrl = append(t.VideoUrl, v.MustGetString(i, "url"))
        }
    }
    t.ExpandedUrl = o.MustGetString("expanded_url")
    return t
}

func (t TwitterEntityUrl) Parse(o *json.JsonObject) TwitterEntityUrl {
    if o == nil {
        return t
    }
    t.DisplayUrl = o.MustGetString("display_url")
    t.ExpandedUrl = o.MustGetString("expanded_url")
    t.Url = o.MustGetString("url")
    return t
}

func (t TwitterEntities) Parse(o *json.JsonObject) TwitterEntities {
    if o == nil {
        return t
    }
    for i := 0; i < o.MustGetArray("hashtags").Length(); i++ {
        t.HashTag = append(t.HashTag, o.MustGetString("hashtags", i, "text"))
    }
    for i := 0; i < o.MustGetArray("media").Length(); i++ {
        t.Media = append(t.Media, TwitterEntityMedia{}.Parse(o.MustGetObject("media", i)))
    }
    for i := 0; i < o.MustGetArray("symbols").Length(); i++ {
        t.Symbols = append(t.Symbols, o.MustGetString("symbols", i, "text"))
    }
    t.Timestamps = nil //TODO
    for i := 0; i < o.MustGetArray("urls").Length(); i++ {
        t.Urls = append(t.Urls, TwitterEntityUrl{}.Parse(o.MustGetObject("urls", i)))
    }
    for i := 0; i < o.MustGetArray("user_mentions").Length(); i++ {
        t.UserMentions = append(t.UserMentions, TwitterEntityUserMention{}.Parse(o.MustGetObject("user_mentions", i)))
    }
    return t
}

func (t Tweet) ParseResult(o *json.JsonObject) Tweet {
    t.RestID = o.MustGetString("rest_id")
    t.CreatedAt, _ = time.Parse(TimeForm, o.MustGetString("legacy", "created_at"))
    t.FullText = o.MustGetString("legacy", "full_text")
    t.Lang = o.MustGetString("legacy", "lang")
    t.PossiblySensitive = o.MustGetBool("legacy", "possibly_sensitive")
    t.Entities = TwitterEntities{}.Parse(o.MustGetObject("legacy", "entities"))
    t.User = TwitterUser{}.Parse(o.MustGetObject("core", "user_results", "result"))
    t.ViewsCount, _ = strconv.ParseInt(o.MustGetString("views", "count"), 0, 64)
    t.BookmarkCount = int(o.MustGetNum("legacy", "bookmark_count"))
    t.FavoriteCount = int(o.MustGetNum("legacy", "favorite_count"))
    t.QuoteCount = int(o.MustGetNum("legacy", "quote_count"))
    t.ReplyCount = int(o.MustGetNum("legacy", "reply_count"))
    t.RetweetCount = int(o.MustGetNum("legacy", "retweet_count"))
    t.Bookmarked = o.MustGetBool("legacy", "bookmarked")
    t.Favorited = o.MustGetBool("legacy", "favorited")
    t.Retweeted = o.MustGetBool("legacy", "retweeted")
    return t
}

func (t Tweet) ParseInstructions(o *json.JsonArray) Tweet {
    for i := 0; i < o.Length(); i++ {
        instruction := o.MustGetObject(i)
        if instruction.MustGetString("type") != "TimelineAddEntries" {
            continue
        }
        entries := instruction.MustGetArray("entries")
        t = Tweet{}.ParseEntry(entries.MustGetObject(0))[0]
        for j := 1; j < entries.Length(); j++ {
            entry := entries.MustGetObject(j)
            if strings.HasPrefix(entry.MustGetString("entryId"), "tweet-") {
                if len(t.ConversationThreads) == 0 {
                    t.ConversationThreads = append(t.ConversationThreads, []Tweet{})
                }
                t.ConversationThreads[0] = append(t.ConversationThreads[0], Tweet{}.ParseEntry(entry)...)
            }
            if strings.HasPrefix(entry.MustGetString("entryId"), "conversationthread-") {
                t.ConversationThreads = append(t.ConversationThreads, Tweet{}.ParseEntry(entry))
            }
        }
    }
    return t
}

func (t Tweet) ParseEntry(o *json.JsonObject) []Tweet {
    if strings.HasPrefix(o.MustGetString("entryId"), "tweet") {
        return []Tweet{Tweet{}.ParseTimelineItem(o.MustGetObject("content"))}
    } else {
        return Tweet{}.ParseTimelineModule(o.MustGetObject("content"))
    }
}

func (t Tweet) ParseTimelineItem(o *json.JsonObject) Tweet {
    item := o.MustGetObject("itemContent")
    if item.MustGetString("itemType") != "TimelineTweet" {
        return Tweet{}
    }
    return Tweet{}.ParseResult(item.MustGetObject("tweet_results", "result"))
}

func (t Tweet) ParseTimelineModule(o *json.JsonObject) []Tweet {
    var tweets []Tweet
    for i := 0; i < o.MustGetArray("items").Length(); i++ {
        item := o.MustGetObject("items", i)
        if res, err := item.GetObject("item", "itemContent", "tweet_results", "result"); err == nil || res != nil {
            tweets = append(tweets, Tweet{}.ParseResult(res))
        }
    }
    return tweets
}
