# reddit irc bot

[![GoDoc](https://godoc.org/github.com/ugjka/ircrssreddit?status.svg)](http://godoc.org/github.com/ugjka/ircrssreddit)
[![Go Report Card](https://goreportcard.com/badge/github.com/ugjka/ircrssreddit)](https://goreportcard.com/report/github.com/ugjka/ircrssreddit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Donate](paypal.svg?raw=true)](https://www.paypal.me/ugjka)

posts newest posts from reddit's rss feeds

## example

```go
package main

import (
    "time"

    bot "github.com/ugjka/ircrssreddit"
)

func main() {
    settings := &bot.Bot{
        Nick:     "examplenick",
        Server:   "chat.freenode.net:6697",
        Channels: []string{"#test"},
        SSL:      true,
        Subreddits: []string{
            "/r/testsub/new/",
            "/r/testsub2/new/",
            "/r/testsub3/new/",
            },
        CheckInterval: time.Minute * 5,
        RoundInterval: time.Minute * 5,
        UserAgent:     "freenode #test personal irc reddit bot",
        PrintSub:      true,
    }
    instance := bot.New(settings)
    instance.Start()
}
```
