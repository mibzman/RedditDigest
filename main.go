package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	RunActions()
}

func RunActions() {
	flag.Parse()
	args := flag.Args()

	var err error
	if args[0] == "serve" {
		err = Serve(args[1])
	} else {
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

func Serve(Filename string) error {
	config, redditBot, err := InitStuff(Filename)
	if err != nil {
		return err
	}

	mx := mux.NewRouter()

	APIManager := APIManager{mx, config, redditBot}
	APIManager.AddRoutes()

	fmt.Println("server is serving")
	fmt.Print(http.ListenAndServe(":"+"8080", mx))
	fmt.Print(http.ListenAndServe(":"+os.Getenv("PORT"), mx))
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
