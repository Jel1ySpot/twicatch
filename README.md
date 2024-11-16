# twicatch

---

> Fetching twitter/X api by playwright.

### Usage

Install:
```shell
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
go get -u github.com/Jel1ySpot/twicatch
```

Example:
```go
package main

import (
    "fmt"
    "github.com/Jel1ySpot/twicatch/pkg/catcher"
)

func main() {
    ctx, err := catcher.CreatePlayWright()
    if err != nil {
        panic(err)
    }
    defer ctx.Close()

    /*
    if err := ctx.LoadCookieFile(cookiePath); err != nil {
        panic(err)
    }
    */

    data, err := ctx.Status("https://x.com/brokenplenty/status/1856147894526718320")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("%v\n", data)
}
```
