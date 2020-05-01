package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterConfig struct {
	ConsumerKey       string `envconfig:"CONSUMER_KEY"`
	ConsumerSecret    string `envconfig:"CONSUMER_SECRET"`
	AccessToken       string `envconfig:"ACCESS_TOKEN"`
	AccessTokenSecret string `envconfig:"ACCESS_TOKEN_SECRET"`
}

// example: go run main.go #MU2020
func main() {
	var c TwitterConfig
	envconfig.Process("TWITTER", &c)
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}

	track := os.Args[1]
	fmt.Println(track)
	filterParams := &twitter.StreamFilterParams{Track: []string{track}}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	stream.Stop()
}
