package main

import (
	"os"
	"time"

	"github.com/turnage/graw/reddit"
)

type RedditBot struct {
	Bot reddit.Bot
}

func InitReddit(data RedditData) (RedditBot, error) {
	file, err := os.Create("reddit.config")
	if err != nil {
		return RedditBot{}, err
	}
	defer file.Close()

	file.WriteString(`user_agent: "` + data.UserAgent + `"
client_id: "` + data.ClientID + `"
client_secret: "` + data.ClientSecret + `"
username: "` + data.Username + `"
password: "` + data.Password + `"`)

	// var err error
	var redditBot RedditBot
	redditBot.Bot, err = reddit.NewBotFromAgentFile("reddit.config", 0)
	return redditBot, err
}

func (bot RedditBot) GetPostsForSub(sub string, limit int, query string, daysOld int, params map[string]string) (Posts, error) {
	harvest, err := bot.Bot.ListingWithParams("/r/"+sub+query, params)
	if err != nil {
		return Posts{[]reddit.Post{}}, err
	}

	var results []reddit.Post

	counter := 0
	for _, post := range harvest.Posts {
		//skips posts if they're older than a day
		if time.Unix(int64(post.CreatedUTC), 0).Before(time.Now().AddDate(0, 0, daysOld*-1)) {
			continue
		}
		results = append(results, *post)
		counter++
		if counter >= limit {
			break
		}
	}
	return Posts{results}, nil
}

type PostsGetter func(Digest) (Posts, error)

func (bot RedditBot) GetMonthlyPosts(Digest Digest) (Posts, error) {
	m := make(map[string]string)
	m["t"] = "month"
	return bot.GetPostsForSub(Digest.Subreddit, Digest.NumPosts, "/top", 30, m)
}

func (bot RedditBot) GetWeeklyPosts(Digest Digest) (Posts, error) {
	m := make(map[string]string)
	m["t"] = "week"
	return bot.GetPostsForSub(Digest.Subreddit, Digest.NumPosts, "/top", 7, m)
}

func (bot RedditBot) GetDailyPosts(Digest Digest) (Posts, error) {
	m := make(map[string]string)
	return bot.GetPostsForSub(Digest.Subreddit, Digest.NumPosts, "", 1, m)
}
