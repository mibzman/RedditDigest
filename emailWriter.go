package main

import (
	"errors"
	"time"
)

type EmailWriter struct {
	Email     string
	redditBot RedditBot
	config    Config
}

func WriteEmail(redditBot RedditBot, config Config) error {
	EmailWriter := EmailWriter{"", redditBot, config}

	EmailWriter.writeHeader()

	err := EmailWriter.writeSecton("Today")
	if err != nil {
		return err
	}

	EmailWriter.writeSpacer()

	if EmailWriter.isWeeklyWeekday() {
		err := EmailWriter.writeSecton("Week")
		if err != nil {
			return err
		}

	}

	EmailWriter.writeSpacer()

	if EmailWriter.isMonthlyDay() {
		err := EmailWriter.writeSecton("Month")
		if err != nil {
			return err
		}
	}

	EmailWriter.writeFooter()

	return EmailWriter.send()

}

func (EmailWriter *EmailWriter) writeSecton(Choice string) error {
	var Result string

	for _, Digest := range EmailWriter.getDigests(Choice) {
		Result += Digest.headline(Choice)

		Digest.populatePosts(EmailWriter.redditBot, Choice)

		FormattedPosts, err := Digest.toString()
		if err != nil {
			return err
		}
		Result += FormattedPosts
	}

	EmailWriter.Email += Result
	return nil
}

func (EmailWriter EmailWriter) getDigests(Choice string) []Digest {
	switch Choice {
	case "Today":
		return EmailWriter.config.DailyDigests
	case "Week":
		return EmailWriter.config.WeeklyDigests
	case "Month":
		return EmailWriter.config.DailyDigests
	}
	panic(errors.New("unknown choice"))
}

func (EmailWriter *EmailWriter) isMonthlyDay() bool {
	return time.Now().Day() == EmailWriter.config.MonthlyDay
}

func (EmailWriter *EmailWriter) isWeeklyWeekday() bool {
	weekday := DayOfTheWeek()
	return weekday == EmailWriter.config.WeeklyWeekday
}

func (EmailWriter *EmailWriter) send() error {
	request := EmailRequest{"", EmailWriter.config.UserEmail, "Reddit Digest for " + DayOfTheWeek(), EmailWriter.Email, []string{}}
	return EmailWriter.config.EmailerConfig.Email(request)
}

func (EmailWriter *EmailWriter) writeHeader() {
	EmailWriter.Email += `<h1>Email Digest for you!</H1> <br></br>`
}

func (EmailWriter *EmailWriter) writeFooter() {
	EmailWriter.Email += "<br></br><br></br>Stay cool <br></br> -RedditDigest Bot"
}

func (EmailWriter *EmailWriter) writeSpacer() {
	EmailWriter.Email += "<br></br><br><hr></hr><hr></hr>"
}
