package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joek/pingu"
)

func main() {
	tux := pingu.NewTux(os.Getenv("TUX_PORT"))
	go tux.Run("5")
	go tux.Run("3")

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	// api.SetLogger(anaconda.BasicLogger)

	f := url.Values{}
	f.Add("track", os.Getenv("TWITTER_TRACK_STRING"))
	stream := api.PublicStreamFilter(f)

	fmt.Println("Twitter connected.")

	for t := range stream.C {
		tux.Wave()
		tweet := t.(anaconda.Tweet)
		fmt.Printf("recieved: %#v\n", tweet.Text)
	}
}
