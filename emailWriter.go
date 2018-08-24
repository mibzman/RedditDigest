package main

import (
	"errors"
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

	err := DigestWriter.writeDigest("Today")
	if err != nil {
		return err
	}

	DigestWriter.writeSpacer()

	if DigestWriter.isWeeklyWeekday() {
		err := DigestWriter.writeDigest("Week")
		if err != nil {
			return err
		}

	}

	DigestWriter.writeSpacer()

	if DigestWriter.isMonthlyDay() {
		err := DigestWriter.writeDigest("Month")
		if err != nil {
			return err
		}
	}

	DigestWriter.writeFooter()

	return DigestWriter.send()

}

func (DigestWriter *DigestWriter) writeDigest(Choice string) error {
	var Result string

	for _, Digest := range DigestWriter.getDigests(Choice) {
		Result += Digest.headline(Choice)

		Digest.populatePosts(DigestWriter.redditBot, Choice)

		FormattedPosts, err := Digest.toString()
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	DigestWriter.Email += Result
	return nil
}

func (DigestWriter DigestWriter) getDigests(Choice string) []Digest {
	switch Choice {
	case "Today":
		return DigestWriter.config.DailyDigests
	case "Week":
		return DigestWriter.config.WeeklyDigests
	case "Month":
		return DigestWriter.config.DailyDigests
	}
	panic(errors.New("unknown choice"))
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

func (Digest Digest) headline(Unit string) string {
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
