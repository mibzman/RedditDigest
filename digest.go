package main

import "fmt"

type Digest struct {
	Subreddit string
	NumPosts  int
	Posts     Posts
}

func (Digest Digest) toString() (string, error) {
	return Digest.Posts.toString()
}

func (digest *Digest) populatePosts(redditBot RedditBot, Choice string) (err error) {
	digest.Posts, err = redditBot.GetPosts(*digest, Choice)

	return
}

func (Digest Digest) headline(Unit string) string {
	return fmt.Sprintf(`<br></br> <h2>This %v's %v Posts from /r/%v </h2>`, Unit, Digest.NumPosts, Digest.Subreddit)
}
