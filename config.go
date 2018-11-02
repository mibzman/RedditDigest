package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Config struct {
	RedditData     RedditData
	UserEmail      string
	EmailerConfig  EmailerConfig
	DailyDigests   []Digest
	WeeklyWeekday  string //The day of the week to send weekly digest
	WeeklyDigests  []Digest
	MonthlyDay     int //The day of the month to send monthly digest
	MonthlyDigests []Digest
}

type EmailerConfig struct {
	URL             string
	Port            int
	SendingAddress  string
	SendingPassword string
}

type RedditData struct {
	UserAgent    string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

func Parse(Filename string) (Config, error) {
	var config Config

	configFile, err := os.Open(Filename)
	if err != nil {
		return Config{}, err
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (data RedditData) CreateAgentFile() error {
	file, err := os.Create("reddit.config")
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(`user_agent: "` + data.UserAgent + `"
client_id: "` + data.ClientID + `"
client_secret: "` + data.ClientSecret + `"
username: "` + data.Username + `"
password: "` + data.Password + `"`)

	return nil
}

func (Config Config) getDigests(Choice string) []Digest {
	switch Choice {
	case "Today":
		return Config.DailyDigests
	case "Week":
		return Config.WeeklyDigests
	case "Month":
		return Config.MonthlyDigests
	}
	panic(errors.New("unknown choice"))
}

func (Config *Config) isMonthlyDay() bool {
	return time.Now().Day() == Config.MonthlyDay
}

func (Config *Config) isWeeklyWeekday() bool {
	weekday := DayOfTheWeek()
	return weekday == Config.WeeklyWeekday
}
