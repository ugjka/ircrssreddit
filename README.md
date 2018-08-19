# reddit irc bot

posts newest posts from reddit's rss feeds

## example

```go
package main

import (
    "time"

    "github.com/ugjka/ircrssreddit"
)

func main() {
    settings := &bot.Bot{
        IrcNick:       "example",
        IrcUser:       "example",
        IrcPass:       "yourpass",
        IrcServer:     "chat.freenode.net:6697",
        IrcChannels:   []string{"#example", "#example2"},
        IrcTLS:        true,
        Endpoints:     []string{"/r/example/new", "/r/example2/new"},
        FetchInterval: time.Minute * 15,
        UserAgent:     "freenode #example irc reddit bot",
    }
    rbot := bot.New(settings)
    rbot.Start()
}
```
