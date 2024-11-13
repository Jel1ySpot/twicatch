package catcher

import (
    cookiemonster "github.com/MercuryEngineering/CookieMonster"
    "github.com/playwright-community/playwright-go"
)

type Context struct {
    PlayWright *playwright.Playwright
    Browser    playwright.Browser
    Cookies    []playwright.OptionalCookie
}

func CreatePlayWright() (*Context, error) {
    pw, err := playwright.Run()
    if err != nil {
        return nil, err
    }
    browser, err := pw.Chromium.Launch()
    if err != nil {
        return nil, err
    }
    return &Context{
        PlayWright: pw,
        Browser:    browser,
    }, nil
}

func (c *Context) AddCookie(cookie playwright.OptionalCookie) {
    c.Cookies = append(c.Cookies, cookie)
}

func (c *Context) LoadCookieFile(path string) error {
    cookies, err := cookiemonster.ParseFile(path)
    if err != nil {
        return err
    }
    c.Cookies = []playwright.OptionalCookie{}
    for _, cookie := range cookies {
        c.Cookies = append(c.Cookies, playwright.Cookie{
            Name:     cookie.Name,
            Value:    cookie.Value,
            Domain:   cookie.Domain,
            Path:     cookie.Path,
            Expires:  float64(cookie.Expires.Unix()),
            HttpOnly: cookie.HttpOnly,
            Secure:   cookie.Secure,
        }.ToOptionalCookie())
    }
    return nil
}

func (c *Context) Close() error {
    if err := c.Browser.Close(); err != nil {
        return err
    }
    return c.PlayWright.Stop()
}
