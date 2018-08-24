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

	DigestWriter.writeSpacer()

	if DigestWriter.isWeeklyWeekday() {
		err := DigestWriter.writeWeeklyDigests()
		if err != nil {
			return err
		}

	}

	DigestWriter.writeSpacer()

	if DigestWriter.isMonthlyDay() {
		err := DigestWriter.writeMonthlyDigests()
		if err != nil {
			return err
		}
	}

	DigestWriter.writeFooter()

	return DigestWriter.send()

}

func (DigestWriter *DigestWriter) writeDigest(postsGetter PostsGetter, Digest Digest, Headline string) (string, error) {
	var Result string

	Posts, err := postsGetter(Digest)
	if err != nil {
		return "", err
	}

	Result += DigestWriter.headline(Headline, Digest)

	FormattedPosts, err := Posts.toString()
	if err != nil {
		return "", err
	}
	Result += FormattedPosts
	return Result, nil
}

func (DigestWriter *DigestWriter) writeDailyDigests() error {
	var Result string
	for _, Digest := range DigestWriter.config.DailyDigests {
		FormattedPosts, err := DigestWriter.writeDigest(DigestWriter.redditBot.GetDailyPosts, Digest, DayOfTheWeek())
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	DigestWriter.Email += Result
	return nil
}

func (DigestWriter *DigestWriter) writeWeeklyDigests() error {
	var Result string

	for _, Digest := range DigestWriter.config.WeeklyDigests {
		FormattedPosts, err := DigestWriter.writeDigest(DigestWriter.redditBot.GetWeeklyPosts, Digest, "Week")
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	DigestWriter.Email += Result
	return nil
}

func (DigestWriter *DigestWriter) writeMonthlyDigests() error {
	var Result string

	for _, Digest := range DigestWriter.config.MonthlyDigests {
		FormattedPosts, err := DigestWriter.writeDigest(DigestWriter.redditBot.GetMonthlyPosts, Digest, "Month")
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	DigestWriter.Email += Result
	return nil
}

func (DigestWriter *DigestWriter) isMonthlyDay() bool {
	return time.Now().Day() == DigestWriter.config.MonthlyDay
}

func (DigestWriter *DigestWriter) isWeeklyWeekday() bool {
	weekday := DayOfTheWeek()
	return weekday == DigestWriter.config.WeeklyWeekday
}

func (DigestWriter *DigestWriter) send() error {
	request := EmailRequest{"", DigestWriter.config.UserEmail, "Reddit Digest for " + DayOfTheWeek(), DigestWriter.Email, []string{}}
	return DigestWriter.config.EmailerConfig.Email(request)
}

func (DigestWriter *DigestWriter) headline(Unit string, Digest Digest) string {
	return fmt.Sprintf(`<br></br> <h2>This %v's %v Posts from /r/%v </h2>`, Unit, Digest.NumPosts, Digest.Subreddit)
}

func (DigestWriter *DigestWriter) writeHeader() {
	DigestWriter.Email += `<h1>Email Digest for you!</H1> <br></br>`
}

func (DigestWriter *DigestWriter) writeFooter() {
	DigestWriter.Email += "<br></br><br></br>Stay cool <br></br> -RedditDigest Bot"
}

func (DigestWriter *DigestWriter) writeSpacer() {
	DigestWriter.Email += "<br></br><br><hr></hr><hr></hr>"
}
