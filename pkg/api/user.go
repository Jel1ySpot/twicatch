package api

import (
    json "github.com/Jel1ySpot/twicatch/pkg/json_helper"
    "time"
)

const (
    UserByScreenNameFormat = `https://x.com/%s`
    UserByRestIDFormat     = `https://x.com/i/user/%s`
)

type (
    TwitterUser struct {
        ID                    string
        Name                  string
        ScreenName            string
        Location              string
        Description           string
        Protected             bool
        Verified              bool
        CreatedAt             time.Time
        Entities              UserEntities
        PinnedTweetIds        []string
        ProfileImageUrlNormal string
        ProfileBannerUrl      string
        FollowersCount        int
        FriendsCount          int
        ListedCount           int
        FavouritesCount       int
        StatusesCount         int
        FollowedBy            *bool
        Following             *bool
        CanDm                 *bool
    }

    UserEntities struct {
        Description TwitterEntities
        Url         TwitterEntities
    }
)

func (u UserEntities) Parse(o *json.JsonObject) UserEntities {
    if o == nil {
        return u
    }
    u.Description = TwitterEntities{}.Parse(o.MustGetObject("description"))
    u.Url = TwitterEntities{}.Parse(o.MustGetObject("url"))
    return u
}

func (t TwitterUser) Parse(o *json.JsonObject) TwitterUser {
    if o == nil {
        return t
    }
    t.ID = o.MustGetString("rest_id")
    t.Name = o.MustGetString("legacy", "name")
    t.ScreenName = o.MustGetString("legacy", "screen_name")
    t.Location = o.MustGetString("legacy", "location")
    t.Description = o.MustGetString("legacy", "description")
    t.Protected = o.MustGetBool("legacy", "protected")
    t.Verified = o.MustGetBool("legacy", "verified")
    t.CreatedAt, _ = time.Parse(TimeForm, o.MustGetString("legacy", "created_at"))
    t.Entities = UserEntities{}.Parse(o.MustGetObject("legacy", "entities"))
    for i := 0; i < o.MustGetArray("legacy", "pinned_tweet_ids_str").Length(); i++ {
        t.PinnedTweetIds = append(t.PinnedTweetIds, o.MustGetString("legacy", "pinned_tweet_ids", i))
    }
    t.ProfileImageUrlNormal = o.MustGetString("legacy", "profile_image_url_https")
    t.ProfileBannerUrl = o.MustGetString("legacy", "profile_banner_url")
    t.FollowersCount = int(o.MustGetNum("legacy", "followers_count"))
    t.FriendsCount = int(o.MustGetNum("legacy", "friends_count"))
    t.ListedCount = int(o.MustGetNum("legacy", "listed_count"))
    t.FavouritesCount = int(o.MustGetNum("legacy", "favourites_count"))
    t.StatusesCount = int(o.MustGetNum("legacy", "statuses_count"))
    if b, err := o.GetBool("legacy", "followed_by"); err == nil {
        t.FollowedBy = &b
    }
    if b, err := o.GetBool("legacy", "following"); err == nil {
        t.Following = &b
    }
    if b, err := o.GetBool("legacy", "can_dm"); err == nil {
        t.CanDm = &b
    }
    return t
}
