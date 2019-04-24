package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

type APIManager struct {
	Router    *mux.Router
	config    Config
	redditBot RedditBot
}

func (APIManager APIManager) AddRoutes() {
	APIManager.Router.HandleFunc("/reddit", APIManager.RedditHandler)
	APIManager.Router.HandleFunc("/redditRSS", APIManager.RedditRSSHandler)
}

func (APIManager APIManager) GetTodaysPosts() []Post {
	Choice := "Today"
	Digests := APIManager.config.getDigests(Choice)

	var Posts []Post
	for _, Digest := range Digests {
		Digest.populatePosts(APIManager.redditBot, Choice)
		Posts = append(Posts, Digest.Posts.toArray()...)
	}
	return Posts
}

func (APIManager APIManager) RedditHandler(w http.ResponseWriter, r *http.Request) {

	Posts := APIManager.GetTodaysPosts()

	encoder := json.NewEncoder(w)
	err := encoder.Encode(Posts)
	if err != nil {
		w.Write([]byte("Failed to encode json:" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (APIManager APIManager) RedditRSSHandler(w http.ResponseWriter, r *http.Request) {

	Posts := APIManager.GetTodaysPosts()
	now := time.Now()

	feed := feeds.Feed{
		Title:       "Reddit RSS",
		Description: "all my reddit stuff",
		Author:      &feeds.Author{Name: "sam borick"},
		Created:     now,
	}

	Items := []*feeds.Item{}

	for _, Post := range Posts {
		Item := &feeds.Item{
			Title:       Post.Title,
			Link:        &feeds.Link{Href: Post.URL},
			Description: Post.SelfTextHTML,
			Created:     now,
		}
		Items = append(Items, Item)
	}

	feed.Items = Items

	rss, err := feed.ToRss()
	if err != nil {
		w.Write([]byte("Failed to encode rss:" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(rss))

}

// func WriteJSONResponse(w http.ResponseWriter, Data ...interface{}) {
// 	encoder := json.NewEncoder(w)
// 	err := encoder.Encode(Data)
// 	if err != nil {
// 		w.Write([]byte("Failed to encode json:" + err.Error()))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// }
