package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/Jel1ySpot/twicatch/pkg/catcher"
)

const HelpText = `**This is a help-message.**

twitake - Catch twitter/x api by using rod.

Usage:
    twitake [args...] command <url>

Example:
    twitake -cookie=.\cookie.txt status https://x.com/Kurarifox/status/1856430975087194466

Note:
    Install deps with
        go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps"`

var cookiePath string

func main() {
    flag.Parse()

    args := flag.Args()

    switch args[0] {
    case "status":
        status(cookiePath, args[1])
    case "user":
        user(cookiePath, args[1])
    default:
        fmt.Println(HelpText)
        return
    }
}

func init() {
    flag.StringVar(&cookiePath, "cookie", "", "cookie.txt")
}

func status(cookiePath string, url string) {
    ctx, err := catcher.GetContext()
    if err != nil {
        panic(err)
    }
    defer ctx.Close()

    if cookiePath != "" {
        if err := ctx.LoadCookieFile(cookiePath); err != nil {
            panic(err)
        }
    }

    data, err := ctx.Status(url)
    if err != nil {
        panic(err)
    }
    //object, err := data.GetObject("data", "threaded_conversation_with_injections_v2", "instructions", 0, "entries", 0, "content", "itemContent", "tweet_results", "result")
    //if err != nil {
    //   panic(err)
    //}
    //s, _ := object.Encode("  ")
    s, _ := json.MarshalIndent(data, "", "  ")
    fmt.Printf("%s\n", s)
}

func user(cookiePath string, url string) {
    ctx, err := catcher.GetContext()
    if err != nil {
        panic(err)
    }
    defer ctx.Close()

    if cookiePath != "" {
        if err := ctx.LoadCookieFile(cookiePath); err != nil {
            panic(err)
        }
    }

    tweets, err := ctx.UserTweets(url)
    if err != nil {
        panic(err)
    }
    //object, err := data.GetObject("data", "threaded_conversation_with_injections_v2", "instructions", 0, "entries", 0, "content", "itemContent", "tweet_results", "result")
    //if err != nil {
    //   panic(err)
    //}
    //s, _ := object.Encode("  ")
    s, _ := json.MarshalIndent(tweets, "", "  ")
    fmt.Printf("%s\n", s)
}
