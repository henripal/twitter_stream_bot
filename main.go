package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

var accessToken = getenv("ACCESS_TOKEN")
var accessSecret = getenv("ACCESS_SECRET")
var consumerKey = getenv("CONSUMER_KEY")
var consumerSecret = getenv("CONSUMER_SECRET")

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("Missing env variable " + name)
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {

	api := anaconda.NewTwitterApiWithCredentials(accessToken,
		accessSecret,
		consumerKey,
		consumerSecret)

	stream := api.PublicStreamFilter(url.Values{
		"track": []string{"NIPS"},
	})

	defer stream.Stop()

	for v := range stream.C {
		t := v.(anaconda.Tweet)
		if (strings.Contains(t.Text, " NIPS ") || strings.Contains(t.Text, " #NIPS ")) && !strings.Contains(t.Text, "#ProtestNIPS") {
			greeting := "Hi, @" + t.User.ScreenName + ". To most people this sounds like: \n"
			editedTweet := strings.Replace(t.Text, "NIPS", "nipples", -1)
			candidateTweet := greeting + editedTweet
			maxLength := min(len(candidateTweet), 280)
			finalTweetWithHashTag := candidateTweet[:maxLength] + "\n #ProtestNIPS"

			fmt.Println("just tweeted: " + finalTweetWithHashTag)

			api.PostTweet(finalTweetWithHashTag, url.Values{
				"in_reply_to_status_id": []string{t.IdStr},
			})
		} else {
			fmt.Println("passed on a vulgar tweet: " + t.Text)
		}
	}
}
