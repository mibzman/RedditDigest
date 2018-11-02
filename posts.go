package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/turnage/graw/reddit"
)

type Posts struct {
	list []Post
}

type Post struct {
	reddit.Post
}

func (Posts *Posts) append(post Post) {
	if len(Posts.list) == 0 {
		Posts.list = []Post{post}
	} else {
		Posts.list = append(Posts.list, post)
	}
}

func (Posts *Posts) MarshalJSON() ([]byte, error) {
	return json.Marshal(Posts.list)
}

func (Posts Posts) toString() (string, error) {
	var Result string
	for _, Post := range Posts.list {
		PostContent, err := Post.toString()
		if err != nil {
			return "", err
		}

		Result += PostContent
	}
	return Result, nil
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

func (Post *Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Title        string
		SelfTextHTML string
		URL          string
	}{
		Title:        Post.Title,
		SelfTextHTML: Post.SelfTextHTML,
		URL:          Post.URL,
	})
}

func (Post Post) isImage() bool {
	FileExtension := Post.URL[len(Post.URL)-3:]
	return FileExtension == "jpg" || FileExtension == "png"
}

func (Post Post) isOlderThan(daysOld int) bool {
	return time.Unix(int64(Post.CreatedUTC), 0).
		Before(time.Now().AddDate(0, 0, daysOld*-1))
}
