package main

import (
	"flag"
	"fmt"
)

func main() {
	RunActions()
}

func RunActions() {
	flag.Parse()
	args := flag.Args()
	err := Run(args[0])
	if err != nil {
		panic(err)
	}
}

func Run(Filename string) error {
	config, err := Parse(Filename)
	if err != nil {
		return err
	}

	redditBot, err := InitReddit(config.RedditData)
	if err != nil {
		return err
	}

	var Email string

	Email += `<h1>Email Digest for you!</H1> <br></br>`

	for _, Digest := range config.DailyDigests {
		posts, err := redditBot.GetPostsForSub(Digest.Subreddit, Digest.NumPosts)
		if err != nil {
			return err
		}

		Email += fmt.Sprintf(`<br></br><hr></hr> <h2>Today's %v Posts from /r/%v </h2>`, Digest.NumPosts, Digest.Subreddit)

		for _, post := range posts {
			if post.IsSelf {
				Email += fmt.Sprintf(`<h3>%v</h3><br></br>`, post.Title)
				Email += fmt.Sprintf(`%v<br></br>`, post.SelfTextHTML)
			} else {
				// Email += fmt.Sprintf("%v", post.IsRedditMediaDomain)
				Email += fmt.Sprintf(`<a href="%v">%v </a><br></br>`, post.URL, post.Title)
			}
		}
	}

	Email += "<br></br><br></br>Stay cool <br></br> -RedditDigest Bot"

	request := EmailRequest{"", config.UserEmail, "Reddit Digest", Email, []string{}}

	return config.EmailerConfig.Email(request)

	// return nil
}
