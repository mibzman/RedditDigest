package main

import "os"

func InitReddit(data RedditData) {
	file, err := os.Create("reddit.config")
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(`user_agent: "` + data.UserAgent + `"
client_id: "` + data.ClientID + `"
client_secret: "` + data.ClientSecret + `"
username: "` + data.Username + `"
password: "` + data.Password + `"`)
}
