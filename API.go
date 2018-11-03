package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type APIManager struct {
	Router    *mux.Router
	config    Config
	redditBot RedditBot
}

func (APIManager APIManager) AddRoutes() {
	APIManager.Router.HandleFunc("/reddit", APIManager.RedditHandler)
}

func (APIManager APIManager) RedditHandler(w http.ResponseWriter, r *http.Request) {
	Choice := "Today"
	Digests := APIManager.config.getDigests(Choice)

	var Posts []Post
	for _, Digest := range Digests {
		Digest.populatePosts(APIManager.redditBot, Choice)
		Posts = append(Posts, Digest.Posts.toArray()...)
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(Posts)
	if err != nil {
		w.Write([]byte("Failed to encode json:" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WriteJSONResponse(w http.ResponseWriter, Data ...interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(Data)
	if err != nil {
		w.Write([]byte("Failed to encode json:" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
