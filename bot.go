package bot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/martinlindhe/base36"
	kitty "github.com/ugjka/kittybot"
	"gopkg.in/inconshreveable/log15.v2"

	"github.com/mmcdole/gofeed"
)

// Client let's you fiddle with http.Client
var Client = &http.Client{}

// Bot contains bot's settings
type Bot struct {
	Server        string
	Nick          string
	Channels      []string
	SSL           bool
	Password      string
	Subreddits    []string
	CheckInterval time.Duration
	RoundInterval time.Duration
	UserAgent     string
	PrintSub      bool
	irc           *kitty.Bot
	highestID     uint64
	send          chan string
	feed          *gofeed.Parser
}

//New creates a new bot object
func New(bot *Bot) *Bot {
	bot.send = make(chan string, 100)
	bot.feed = gofeed.NewParser()
	bot.irc = kitty.NewBot(bot.Server, bot.Nick, func(irc *kitty.Bot) {
		irc.Channels = bot.Channels
		irc.SSL = bot.SSL
		irc.Password = bot.Password
	})
	return bot
}

func (bot *Bot) printer() {
	irc := bot.irc
	for msg := range bot.send {
		for _, ch := range bot.Channels {
			irc.Msg(ch, msg)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (bot *Bot) ircLoop() {
	irc := bot.irc
	logHandler := log15.LvlFilterHandler(log15.LvlInfo, log15.StdoutHandler)
	irc.Logger.SetHandler(logHandler)
	for {
		irc.Run()
		irc.Info("reconnecting...")
		time.Sleep(time.Second * 30)
	}
}

// Get posts
func (bot *Bot) fetch(endpoint string) (p *gofeed.Feed, err error) {
	if !strings.Contains(endpoint, ".rss") {
		endpoint += ".rss"
	}
	req, err := http.NewRequest("GET", "https://www.reddit.com"+endpoint, nil)
	if err != nil {
		return nil, err
	}
	// Headers.
	req.Header.Set("User-Agent", bot.UserAgent)

	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fetch response error: " + resp.Status)
	}
	return bot.feed.Parse(resp.Body)
}

func (bot *Bot) firstRun() error {
	for _, v := range bot.Subreddits {
		posts, err := bot.fetch(v)
		if err != nil {
			log.Println("first run: ", err)
			return err
		}
		for _, v := range posts.Items {
			if !strings.HasPrefix(v.GUID, "t3_") {
				continue
			}
			decoded := base36.Decode(v.GUID[3:])
			if bot.highestID < decoded {
				bot.highestID = decoded
			}
		}
	}
	return nil
}

func (bot *Bot) fetchPosts() {
	reddit := "reddit"
	var high uint64
	store := make(map[uint64]bool)
	for _, sub := range bot.Subreddits {
		posts, err := bot.fetch(sub)
		if err != nil {
			log.Println("could not fetch posts: ", err)
			return
		}
		for _, item := range posts.Items {
			if !strings.HasPrefix(item.GUID, "t3_") {
				continue
			}
			decoded := base36.Decode(item.GUID[3:])
			if _, ok := store[decoded]; ok {
				continue
			}
			store[decoded] = true
			if high < decoded {
				high = decoded
			}
			if bot.highestID < decoded {
				name := ""
				if item.Author == nil {
					name = "account_deleted"
				} else {
					name = item.Author.Name
				}
				if bot.PrintSub && item.Categories != nil {
					reddit = "/r/" + item.Categories[0]
				}
				bot.send <- fmt.Sprintf("[%s] [%s] %s https://redd.it/%s", reddit, name, item.Title, item.GUID[3:])
			}
		}
	}
	bot.highestID = high
}

//Start starts the bot
func (bot *Bot) Start() {
	for {
		err := bot.firstRun()
		if err == nil {
			log.Println("first run succeeded")
			break
		}
		log.Println("first run failed: ", err)
		time.Sleep(time.Minute * 10)
		log.Println("retrying first run")
	}
	go bot.printer()
	go bot.ircLoop()

	rounded := time.Now().Round(bot.RoundInterval)
	if time.Now().After(rounded) {
		rounded = rounded.Add(bot.RoundInterval)
	}
	time.Sleep(rounded.Sub(time.Now()))
	ticker := time.NewTicker(bot.CheckInterval)
	bot.fetchPosts()
	for {
		select {
		case <-ticker.C:
			bot.fetchPosts()
		}
	}
}
