package main

import (
	"time"

	"github.com/turnage/graw/reddit"
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

func (DigestWriter DigestWriter) send() error {
	request := EmailRequest{"", DigestWriter.config.UserEmail, "Reddit Digest for " + DayOfTheWeek(), DigestWriter.Email, []string{}}
	return DigestWriter.config.EmailerConfig.Email(request)
}

func (DigestWriter DigestWriter) writeDailyDigests() error {
	// var Result string
	// for _, Digest := range DigestWriter.config.DailyDigests {
	// 	Posts, err := DigestWriter.redditBot.GetDailyPosts(Digest)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	Result += DigestWriter.headline("Today", Digest)

	// 	FormattedPosts, err := DigestWriter.postsToString(Posts)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	Result += FormattedPosts
	// }

	// DigestWriter.Email += Result
	// return nil
	return DigestWriter.writeDigests(DigestWriter.redditBot.GetDailyPosts, "Today")
}

func (DigestWriter DigestWriter) writeWeeklyDigests() error {
	// var Result string

	// for _, Digest := range DigestWriter.config.WeeklyDigests {
	// 	Posts, err := DigestWriter.redditBot.GetWeeklyPosts(Digest)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	Result += DigestWriter.headline("Week", Digest)

	// 	FormattedPosts, err := DigestWriter.postsToString(Posts)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	Result += FormattedPosts
	// }

	// DigestWriter.Email += Result
	// return nil
	return DigestWriter.writeDigests(DigestWriter.redditBot.GetWeeklyPosts, "Week")
}

func (DigestWriter DigestWriter) writeMonthlyDigests() error {
	// var Result string

	// for _, Digest := range DigestWriter.config.MonthlyDigests {
	// 	Posts, err := DigestWriter.redditBot.GetMonthlyPosts(Digest)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	Result += DigestWriter.headline("Month", Digest)

	// 	FormattedPosts, err := DigestWriter.postsToString(Posts)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	Result += FormattedPosts
	// }

	// DigestWriter.Email += Result
	// return nil

	return DigestWriter.writeDigests(DigestWriter.redditBot.GetMonthlyPosts, "Month")
}

type PostGetter func(Digest) ([]reddit.Post, error)

func (DigestWriter DigestWriter) writeDigests(postGetter PostGetter, Headline string) error {
	var Result string

	for _, Digest := range DigestWriter.config.MonthlyDigests {
		FormattedDigest, err := DigestWriter.formatDigest(postGetter, Headline, Digest)
		if err != nil {
			return err
		}
		Result += FormattedDigest
	}

	DigestWriter.Email += Result
	return nil
}

func (DigestWriter DigestWriter) formatDigest(postGetter PostGetter, Headline string, Digest Digest) (string, error) {
	Result := DigestWriter.headline(Headline, Digest)

	Posts, err := postGetter(Digest)
	if err != nil {
		return "", err
	}

	FormattedPosts, err := DigestWriter.postsToString(Posts)
	if err != nil {
		return "", err
	}
	Result += FormattedPosts
	return Result, nil

}

func (DigestWriter DigestWriter) isMonthlyDay() bool {
	return time.Now().Day() == DigestWriter.config.MonthlyDay
}

func (DigestWriter DigestWriter) isWeeklyWeekday() bool {
	weekday := DayOfTheWeek()
	return weekday == DigestWriter.config.WeeklyWeekday
}
