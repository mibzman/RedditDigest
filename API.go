package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	APIManager.Router.HandleFunc("/example", APIManager.ExampleRSSHandler)
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

	fmt.Println(Posts)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(Posts)
	if err != nil {
		w.Write([]byte("Failed to encode json:" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (APIManager APIManager) RedditRSSHandler(w http.ResponseWriter, r *http.Request) {

	var Posts []Post

	jsonFile, err := os.Open("./dump")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Posts)

	fmt.Println(Posts)

	feed := &feeds.Feed{
		Title:       "Reddit RSS",
		Link:        &feeds.Link{Href: "https://upaper.herokuapp.com/redditRSS"},
		Description: "all my reddit stuff",
		Author:      &feeds.Author{Name: "Sam Borick", Email: "sam@borick.net"},
	}

	var Items []*feeds.Item

	for _, Post := range Posts {
		tm := time.Unix(int64(Post.CreatedUTC), 0)

		Image := ""

		if Post.isImage() {
			Image = `<img src="` + Post.URL + `" />`
		}

		Description := Image + Post.SelfTextHTML

		Item := &feeds.Item{
			Title:       Post.Title,
			Link:        &feeds.Link{Href: Post.URL},
			Id:          Post.ID,
			Description: Description,
			Created:     tm,
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

	// fmt.Println(rss)

	w.Write([]byte(rss))
	// w.Write([]byte("hi"))

}

func (APIManager APIManager) ExampleRSSHandler(w http.ResponseWriter, r *http.Request) {

	now := time.Now()
	// feed := &feeds.Feed{
	// 	Title:       "jmoiron.net blog",
	// 	Link:        &feeds.Link{Href: "http://jmoiron.net/blog"},
	// 	Description: "discussion about tech, footie, photos",
	// 	Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
	// 	Created:     now,
	// }
	feed := &feeds.Feed{
		Title:       "Reddit RSS",
		Link:        &feeds.Link{Href: "http://jmoiron.net/blog"},
		Description: "all my reddit stuff",
		Author:      &feeds.Author{Name: "sam borick"},
	}

	feed.Items = []*feeds.Item{
		&feeds.Item{
			Title:       "Limiting Concurrency in Go",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Description: "A discussion on controlled parallelism in golang",
			Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     now,
		},
		&feeds.Item{
			Title:       "Logic-less Template Redux",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
			Description: "More thoughts on logicless templates",
			Created:     now,
		},
		&feeds.Item{
			Title:       "Idiomatic Code Reuse in Go",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
			Description: "How to use interfaces <em>effectively</em>",
			Created:     now,
		},
	}

	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	json, err := feed.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(atom, "\n", rss, "\n", json)
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
