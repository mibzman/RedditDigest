package main

import (
	"github.com/turnage/graw/reddit"
)

type RedditBot struct {
	Bot reddit.Bot
}

func InitReddit(data RedditData) (RedditBot, error) {
	err := data.CreateAgentFile()
	if err != nil {
		return RedditBot{}, err
	}

	var redditBot RedditBot
	redditBot.Bot, err = reddit.NewBotFromAgentFile("reddit.config", 0)
	return redditBot, err
}

func (bot RedditBot) GetPosts(Digest Digest, choice string) (Posts, error) {
	query, daysOld, params := PickParams(choice)
	harvest, err := bot.Bot.ListingWithParams("/r/"+Digest.Subreddit+query, params)
	if err != nil {
		return Posts{}, err
	}

	var results Posts

	counter := 0
	for _, post := range harvest.Posts {
		Post := Post{*post}

		if Post.isOlderThan(daysOld) {
			continue
		}
		results.append(Post)
		counter++
		if counter >= Digest.NumPosts {
			break
		}
	}
	return results, nil
}

func PickParams(choice string) (query string, daysOld int, params map[string]string) {
	params = make(map[string]string)
	switch choice {
	case "Today":
		query = ""
		daysOld = 1
	case "Week":
		query = "/top"
		daysOld = 7
		params["t"] = "week"
	case "Month":
		query = "/top"
		daysOld = 30
		params["t"] = "month"
	}
	return
}
