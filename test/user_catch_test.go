package test

import (
    "fmt"
    "github.com/Jel1ySpot/twicatch/pkg/api"
    "github.com/Jel1ySpot/twicatch/pkg/catcher"
    "os"
    "testing"
)

func TestUserCatch(t *testing.T) {
    ctx, err := catcher.CreatePlayWright()
    if err != nil {
        t.Fatal(err)
    }

    cookiePath := os.Getenv("cookie")
    if cookiePath != "" {
        t.Log("loading cookie file")
        if err := ctx.LoadCookieFile(cookiePath); err != nil {
            t.Fatal(err)
        }
    }

    tweets, err := ctx.UserTweets(fmt.Sprintf(api.UserByScreenNameFormat, "uiokv_829"))
    if err != nil {
        t.Fatal(err)
    }
    if tweets == nil {
        t.Fatal("fetch data is nil")
    }
}
