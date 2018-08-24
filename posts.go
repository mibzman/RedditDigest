package main

import (
	"fmt"

	"github.com/turnage/graw/reddit"
)

type Posts struct {
	list []reddit.Post
}

type Post struct {
	reddit.Post
}

func (Posts Posts) toString() (string, error) {
	var Result string
	for _, post := range Posts.list {
		Post := Post{post}
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

func (Post Post) isImage() bool {
	FileExtension := Post.URL[len(Post.URL)-3:]
	return FileExtension == "jpg" || FileExtension == "png"
}
