package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jasonlvhit/gocron"
)

var HasRun bool

func main() {
	HasRun = false
	gocron.Every(1).Day().At("13:00").Do(RunMibzman)

	for {
		if HasRun {
			break
		}

		time.Sleep(time.Minute)
	}
}

func RunMibzman() {
	err := Run("mibzman.config")
	if err != nil {
		panic(err)
	}
	log.Println("sent email")
	HasRun = true
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
		posts, err := redditBot.GetDailyPostsForSub(Digest.Subreddit, Digest.NumPosts)
		if err != nil {
			return err
		}

		Email += fmt.Sprintf(`<br></br><hr></hr> <h2>Today's %v Posts from /r/%v </h2>`, Digest.NumPosts, Digest.Subreddit)

		for _, post := range posts {
			result, err := GeneratePostEmailContent(post)
			if err != nil {
				return err
			}

			Email += result
		}
	}

	weekday := time.Now().Weekday().String()

	if weekday == config.WeeklyWeekday {
		Email += "<br></br><br><hr></hr><hr></hr>"

		for _, Digest := range config.WeeklyDigests {
			posts, err := redditBot.GetWeeklyPostsForSub(Digest.Subreddit, Digest.NumPosts)
			if err != nil {
				return err
			}

			Email += fmt.Sprintf(`<br></br> <h2>This Week's %v Posts from /r/%v </h2>`, Digest.NumPosts, Digest.Subreddit)

			for _, post := range posts {
				result, err := GeneratePostEmailContent(post)
				if err != nil {
					return err
				}

				Email += result
			}
		}

	}

	if time.Now().Day() == config.MonthlyDay {
		Email += "<br></br><br><hr></hr><hr></hr>"

		for _, Digest := range config.MonthlyDigests {
			posts, err := redditBot.GetMonthlyPostsForSub(Digest.Subreddit, Digest.NumPosts)
			if err != nil {
				return err
			}

			Email += fmt.Sprintf(`<br></br> <h2>This Months's %v Posts from /r/%v </h2>`, Digest.NumPosts, Digest.Subreddit)

			for _, post := range posts {
				result, err := GeneratePostEmailContent(post)
				if err != nil {
					return err
				}

				Email += result
			}
		}
	}

	Email += "<br></br><br></br>Stay cool <br></br> -RedditDigest Bot"

	request := EmailRequest{"", config.UserEmail, "Reddit Digest for " + weekday, Email, []string{}}

	return config.EmailerConfig.Email(request)

	// return nil
}
