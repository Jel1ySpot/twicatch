package test

import (
    "github.com/Jel1ySpot/twicatch/pkg/catcher"
    "os"
    "testing"
)

func TestStatusCatch(t *testing.T) {
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

    data, err := ctx.Status("https://x.com/lingquantang/status/1856270788632162581")
    if err != nil {
        t.Fatal(err)
    }
    if data == nil {
        t.Fatal("fetch data is nil")
    }
}
