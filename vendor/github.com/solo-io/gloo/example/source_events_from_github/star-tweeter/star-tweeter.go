package main

import (
	"fmt"

	"os"

	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ilackarms/go-github-webhook-server/github"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/example/source_events_from_github/base"
	"github.com/solo-io/gloo/pkg/log"
)

func main() {
	log.Fatalf("err:", run())
}

func run() error {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	log.Printf("starting twitter client")

	opts := base.Opts{
		ClientID:  os.Getenv("HOSTNAME"),
		ClusterID: "test-cluster",
		NatsURL:   "nats://nats-streaming.default.svc.cluster.local:4222",
		Subject:   "github-webhooks",
		Handler:   handleWatch(client),
	}
	base.Run(opts)
	log.Printf("terminated")
	return nil
}

func handleWatch(client *twitter.Client) func(watch *github.WatchEvent) error {
	return func(watch *github.WatchEvent) error {
		sender := watch.Sender.Login
		text := fmt.Sprintf("Thanks to %v for the star on github.com/solo-io/gloo (%v)!", sender, time.Now().Format(time.RFC850))
		// Send a Tweet
		tweet, resp, err := client.Statuses.Update(text, nil)
		if err != nil {
			return errors.Wrapf(err, "sending tweet, resp: %v", resp)
		}
		log.Printf("successful tweet: ", tweet)

		return nil
	}
}
