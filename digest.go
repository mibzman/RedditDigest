package main

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
