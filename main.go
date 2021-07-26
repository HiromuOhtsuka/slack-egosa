package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/slack-go/slack"
)

const (
	// Required
	SlackToken = "SLACK_TOKEN"
	// Required
	WebhookURL = "WEBHOOK_URL"
	// Required
	// , 区切りで与える
	// 例: hoge,fuga
	Keywords = "KEYWORDS"
	// Default: 20
	// 最大 100 まで
	MaxSearchCount = "MAX_SEARCH_COUNT"
	// Default: 24
	DurationHours = "DURATION_HOURS"
	// 検索から除外するチャンネルのリスト
	ExcludeChannels = "EXCLUDE_CHANNELS"
	// 検索から除外するユーザのリスト
	ExcludeUsers = "EXCLUDE_USERS"
	Debug        = "DEBUG"
)

type Config struct {
	SlackToken      string
	WebhookURL      string
	Keywords        []string
	MaxSearchCount  int
	DurationHours   int
	ExcludeChannels []string
	ExcludeUsers    []string
	Debug           bool
}

type Message struct {
	Keyword   string
	Channel   string
	Permalink string
}

func (m Message) String() string {
	return fmt.Sprintf(":eyes: キーワード %s に関する発言が #%s でありました\n%s", m.Keyword, m.Channel, m.Permalink)
}

func main() {
	config := readEnv()
	threshold := time.Now().Add(-time.Duration(3600*config.DurationHours) * time.Second)
	older := func(sm slack.SearchMessage) bool {
		timestamp, err := parseTimestamp(sm.Timestamp)
		if err != nil {
			panic(err)
		}
		return timestamp.Before(threshold)
	}
	api := slack.New(config.SlackToken)

	for _, keyword := range config.Keywords {
		slackMessages, err := api.SearchMessages(keyword, slack.SearchParameters{
			Sort:          "timestamp",
			SortDirection: "desc",
			Highlight:     false,
			Count:         config.MaxSearchCount,
			Page:          1,
		})
		if err != nil {
			panic(err)
		}
		for _, sm := range slackMessages.Matches {
			if older(sm) {
				continue
			}
			if containsArray(sm.Channel.Name, config.ExcludeChannels) {
				continue
			}
			if len(sm.User) == 0 {
				continue
			}
			if containsArray(sm.Username, config.ExcludeUsers) {
				continue
			}
			userinfo, err := api.GetUserInfo(sm.User)
			if err != nil {
				panic(err)
			}
			if userinfo.IsBot || userinfo.IsAppUser {
				continue
			}
			message := Message{
				Keyword:   keyword,
				Channel:   sm.Channel.Name,
				Permalink: sm.Permalink,
			}
			if !config.Debug {
				if err := postMessage(config.WebhookURL, message.String()); err != nil {
					panic(err)
				}
			} else {
				fmt.Println(message)
			}
		}
	}
}

func readEnv() Config {
	slackToken := os.Getenv(SlackToken)
	if len(slackToken) == 0 {
		panic(SlackToken + " is empty. must be not empty.")
	}
	webhookURL := os.Getenv(WebhookURL)
	if len(webhookURL) == 0 {
		panic(WebhookURL + " is empty. must be not empty.")
	}
	keywords := strings.Split(os.Getenv(Keywords), ",")
	if len(keywords) == 0 {
		panic(Keywords + " is empty. must be not empty.")
	}
	maxSearchCount := 20
	if len(os.Getenv(MaxSearchCount)) != 0 {
		value, err := strconv.Atoi(os.Getenv(MaxSearchCount))
		if err != nil {
			panic(err)
		}
		maxSearchCount = value
	}
	durationHours := 24
	if len(os.Getenv(DurationHours)) != 0 {
		value, err := strconv.Atoi(os.Getenv(DurationHours))
		if err != nil {
			panic(err)
		}
		durationHours = value
	}
	excludeChannels := strings.Split(os.Getenv(ExcludeChannels), ",")
	excludeUsers := strings.Split(os.Getenv(ExcludeUsers), ",")
	debug := false
	if len(os.Getenv(Debug)) != 0 {
		debug = true
	}
	return Config{
		SlackToken:      slackToken,
		WebhookURL:      webhookURL,
		Keywords:        keywords,
		MaxSearchCount:  maxSearchCount,
		DurationHours:   durationHours,
		ExcludeChannels: excludeChannels,
		ExcludeUsers:    excludeUsers,
		Debug:           debug,
	}
}

func parseTimestamp(timestamp string) (time.Time, error) {
	value, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		return time.Time{}, err
	}
	// 小数点は切り捨ててしまう
	return time.Unix(int64(value), 0), nil
}

func postMessage(webhookURL string, message string) error {
	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]interface{}{"text": message, "unfurl_links": true}).
		Post(webhookURL)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("StatusCode = %d. not 200.", resp.StatusCode())
	}
	return nil
}

func containsArray(element string, array []string) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}