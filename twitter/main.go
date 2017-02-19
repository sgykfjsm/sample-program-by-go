package main

import (
	"fmt"
	"log"

	"net/url"

	"github.com/BurntSushi/toml"
	"github.com/ChimeraCoder/anaconda"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	credential Credential
	api        *anaconda.TwitterApi
)

type Credential struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type SavedTweet struct {
	TweetID    int64
	ScreenName string
	Name       string
}

func (s *SavedTweet) Save(id int64, screenName, name string) {
	s.TweetID = id
	s.ScreenName = screenName
	s.Name = name
}

func init() {
	tomlFile := "credential.toml"
	if _, err := toml.DecodeFile(tomlFile, &credential); err != nil {
		log.Fatalf("Failed to read %s: %s", tomlFile, err.Error())
	}

	anaconda.SetConsumerKey(credential.ConsumerKey)
	anaconda.SetConsumerSecret(credential.ConsumerSecret)

	api = anaconda.NewTwitterApi(credential.AccessToken, credential.AccessTokenSecret)
}

func main() {
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	fmt.Printf("Consumer Key: %s\n", credential.ConsumerKey)
	fmt.Printf("Consumer Secret: %s\n", credential.ConsumerSecret)

	println()

	var savedTweet SavedTweet
	v := url.Values{}
	v.Set("count", "30")
	v.Set("result_type", "recent")
	v.Set("src", "typd")

	go func() {
		t := time.NewTicker(10 * time.Second)
		t2 := time.NewTicker(120 * time.Second)

		tweetCount := 0
		for {
			searchResult, err := api.GetSearch("from:dog_fooder", v)
			if err != nil {
				log.Fatal("Failed to get search result: %s", err.Error())
			}
			showAll := false
			select {
			case <-t2.C:
				tweetCount += 1
				api.PostTweet(fmt.Sprintf("これ%d回目のツイートでした", tweetCount), nil)
			case <-t.C:
				for i, tweet := range searchResult.Statuses {
					fmt.Printf("SavedTweetID: %d\n", savedTweet.TweetID)
					if savedTweet.TweetID == 0 {
						showAll = true
					}
					if !showAll && tweet.Id <= savedTweet.TweetID {
						log.Println("New Tweet Not Found")
						break
					}
					if i == 0 { // Responded Tweets should be sorted from newest to oldest
						savedTweet.Save(tweet.Id, tweet.User.ScreenName, tweet.User.Name)
						log.Println("Saved Tweet")
					}
					log.Printf("%02d:%d:%s(%s): %s at %s\n",
						i, tweet.Id, tweet.User.ScreenName, tweet.User.Name, tweet.Text, tweet.CreatedAt)
				}
				showAll = false
				log.Println("Waiting for 5 seconds...")
			case <-sig:
				done <- true
				t.Stop()
				t2.Stop()
				return
			}
		}
	}()

	<-done
	println()

	pages := api.GetFollowersListAll(nil)
	for page := range pages {
		for i, follower := range page.Followers {
			fmt.Printf("%02d: %s(%s): %s\n", i, follower.ScreenName, follower.Name, follower.Description)
		}
	}
}
