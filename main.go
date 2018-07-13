package main

import (
	"flag"
)

func main() {
	RunActions()
}

func RunActions() {
	flag.Parse()
	args := flag.Args()
	err := Run(args[0])
	if err != nil {
		panic(err)
	}
}

func Run(Filename string) error {
	config, err := Parse(Filename)
	if err != nil {
		return err
	}

	InitReddit(config.RedditData)
	return nil
}
