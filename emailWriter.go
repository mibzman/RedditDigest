package main

import (
	"fmt"
	"time"
)

type DigestWriter struct {
	Email     string
	redditBot RedditBot
	config    Config
}

func WriteEmail(redditBot RedditBot, config Config) error {
	DigestWriter := DigestWriter{"", redditBot, config}

	DigestWriter.writeHeader()

	err := DigestWriter.writeDailyDigests()
	if err != nil {
		return err
	}

	if DigestWriter.isWeeklyWeekday() {
		err := DigestWriter.writeWeeklyDigests()
		if err != nil {
			return err
		}

	}

	if DigestWriter.isMonthlyDay() {
		err := DigestWriter.writeMonthlyDigests()
		if err != nil {
			return err
		}
	}

	DigestWriter.writeFooter()

	return DigestWriter.send()

}

func (DigestWriter DigestWriter) writeHeader() {
	DigestWriter.Email += `<h1>Email Digest for you!</H1> <br></br>`
}

func (DigestWriter DigestWriter) writeFooter() {
	DigestWriter.Email += "<br></br><br></br>Stay cool <br></br> -RedditDigest Bot"
}

func (DigestWriter DigestWriter) send() error {
	request := EmailRequest{"", DigestWriter.config.UserEmail, "Reddit Digest for " + DayOfTheWeek(), DigestWriter.Email, []string{}}
	return DigestWriter.config.EmailerConfig.Email(request)
}

func (DigestWriter DigestWriter) writeDailyDigests() error {
	var Result string
	for _, Digest := range DigestWriter.config.DailyDigests {
		Posts, err := DigestWriter.redditBot.GetDailyPostsForSub(Digest.Subreddit, Digest.NumPosts)
		if err != nil {
			return err
		}

		Result += fmt.Sprintf(`<br></br><hr></hr> <h2>Today's %v Posts from /r/%v </h2>`, Digest.NumPosts, Digest.Subreddit)

		FormattedPosts, err := DigestWriter.postsToString(Posts)
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	DigestWriter.Email += Result
	return nil
}

func (DigestWriter DigestWriter) isWeeklyWeekday() bool {
	weekday := DayOfTheWeek()
	return weekday == DigestWriter.config.WeeklyWeekday
}

func (DigestWriter DigestWriter) writeWeeklyDigests() error {
	var Result string
	Result += "<br></br><br><hr></hr><hr></hr>"

	for _, Digest := range DigestWriter.config.WeeklyDigests {
		Posts, err := DigestWriter.redditBot.GetWeeklyPostsForSub(Digest.Subreddit, Digest.NumPosts)
		if err != nil {
			return err
		}

		Result += fmt.Sprintf(`<br></br> <h2>This Week's %v Posts from /r/%v </h2>`, Digest.NumPosts, Digest.Subreddit)

		FormattedPosts, err := DigestWriter.postsToString(Posts)
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	DigestWriter.Email += Result
	return nil
}

func (DigestWriter DigestWriter) isMonthlyDay() bool {
	return time.Now().Day() == DigestWriter.config.MonthlyDay
}

func (DigestWriter DigestWriter) writeMonthlyDigests() error {
	var Result string
	Result += "<br></br><br><hr></hr><hr></hr>"

	for _, Digest := range DigestWriter.config.MonthlyDigests {
		Posts, err := DigestWriter.redditBot.GetMonthlyPostsForSub(Digest.Subreddit, Digest.NumPosts)
		if err != nil {
			return err
		}

		Result += fmt.Sprintf(`<br></br> <h2>This Months's %v Posts from /r/%v </h2>`, Digest.NumPosts, Digest.Subreddit)

		FormattedPosts, err := DigestWriter.postsToString(Posts)
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	DigestWriter.Email += Result
	return nil
}
