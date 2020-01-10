package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	RunActions()
}

func RunActions() {
	flag.Parse()
	args := flag.Args()

	var err error
	if args[0] == "" {

	} else {

	}

	switch args[0] {
	case "serve":
		Dump(args[1])
		err = Serve(args[1])
	case "dump":
		err = Dump(args[1])
	default:
		err = Run(args[0])
	}

	if err != nil {
		panic(err)
	}
}

func Run(Filename string) error {
	config, redditBot, err := InitStuff(Filename)
	if err != nil {
		return err
	}
	return WriteEmail(redditBot, config)
}

func Dump(Filename string) error {
	config, redditBot, err := InitStuff(Filename)
	if err != nil {
		return err
	}
	Posts := GetTodaysPostsMain(redditBot, config)

	json, err := json.Marshal(Posts)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return ioutil.WriteFile("./dump", []byte(json), 0644)
}

func GetTodaysPostsMain(redditBot RedditBot, config Config) []Post {
	Choice := "Today"
	Digests := config.getDigests(Choice)

	var Posts []Post
	for _, Digest := range Digests {
		Digest.populatePosts(redditBot, Choice)
		Posts = append(Posts, Digest.Posts.toArray()...)
	}
	return Posts
}

func Serve(Filename string) error {
	config, redditBot, err := InitStuff(Filename)
	if err != nil {
		return err
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // "http://localhost:3000", "http://localhost:4200"
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"}, //Authorization
	})

	mx := mux.NewRouter()

	APIManager := APIManager{mx, config, redditBot}
	APIManager.AddRoutes()

	handler := c.Handler(mx)

	port := os.Getenv("PORT")

	fmt.Println("server is serving on ", port)
	// fmt.Print(http.ListenAndServe(":"+"8081", handler))
	fmt.Print(http.ListenAndServe(":"+port, handler))
	return nil
}

func InitStuff(Filename string) (Config, RedditBot, error) {
	config, err := Parse(Filename)
	if err != nil {
		return Config{}, RedditBot{}, err
	}

	redditBot, err := InitReddit(config.RedditData)
	if err != nil {
		return Config{}, RedditBot{}, err
	}

	return config, redditBot, nil
}
