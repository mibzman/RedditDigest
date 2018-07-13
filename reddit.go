package main

import (
	"os"
	"time"

	"github.com/turnage/graw/reddit"
)

type RedditBot struct {
	Bot reddit.Bot
}

// type Post struct {
// 	Title string
// 	Link string
// }

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

func (bot RedditBot) GetPostsForSub(sub string, limit int) ([]reddit.Post, error) {
	harvest, err := bot.Bot.Listing("/r/"+sub, "")
	if err != nil {
		return []reddit.Post{}, err
	}

	var results []reddit.Post

	counter := 0
	for _, post := range harvest.Posts[:limit*2] {
		//skips posts if they're older than a day
		if time.Unix(int64(post.CreatedUTC), 0).Before(time.Now().AddDate(0, 0, -1)) {
			continue
		}
		results = append(results, *post)
		counter++
		if counter >= limit {
			break
		}
		// fmt.Printf("[%s] posted [%s]\n", post.Author, post.Title)
	}
	return results, nil
}
