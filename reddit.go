package main

import (
	"fmt"
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

func (bot RedditBot) GetPostsForSub(sub string, limit int, query string, daysOld int, params map[string]string) ([]reddit.Post, error) {
	harvest, err := bot.Bot.ListingWithParams("/r/"+sub+query, params)
	if err != nil {
		return []reddit.Post{}, err
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
	return results, nil
}

type PostsGetter func(string, int) ([]reddit.Post, error)

func (bot RedditBot) GetMonthlyPosts(Digest Digest) ([]reddit.Post, error) {
	m := make(map[string]string)
	m["t"] = "month"
	return bot.GetPostsForSub(Digest.Subreddit, Digest.NumPosts, "/top", 30, m)
}

func (bot RedditBot) GetWeeklyPosts(Digest Digest) ([]reddit.Post, error) {
	m := make(map[string]string)
	m["t"] = "week"
	return bot.GetPostsForSub(Digest.Subreddit, Digest.NumPosts, "/top", 7, m)
}

func (bot RedditBot) GetDailyPosts(Digest Digest) ([]reddit.Post, error) {
	m := make(map[string]string)
	return bot.GetPostsForSub(Digest.Subreddit, Digest.NumPosts, "", 1, m)
}

func (DigestWriter DigestWriter) postsToString(Posts []reddit.Post) (string, error) {
	var Result string
	for _, post := range Posts {
		Post := Post{post}
		PostContent, err := Post.toString()
		if err != nil {
			return "", err
		}

		Result += PostContent
	}
	return Result, nil
}

type Post struct {
	reddit.Post
}

func (Post Post) toString() (result string, err error) {
	if Post.IsSelf { //is a self post
		result += fmt.Sprintf(`<h3>%v</h3>`, Post.Title)
		result += fmt.Sprintf(`%v<br></br>`, Post.SelfTextHTML)
	} else {
		if Post.isImage() {
			result += fmt.Sprintf(`<h4>%v </h4> <img src="%v" width="500"> </img> <br></br><br></br>`, Post.Title, Post.URL)
		} else {
			result += fmt.Sprintf(`<a href="%v">%v </a> <br></br><br></br>`, Post.URL, Post.Title)
		}

	}
	return
}

func (Post Post) isImage() bool {
	FileExtension := Post.URL[len(Post.URL)-3:]
	return FileExtension == "jpg" || FileExtension == "png"
}
