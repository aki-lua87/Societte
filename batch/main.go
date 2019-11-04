package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	api           *anaconda.TwitterApi
	exeCounter    int
	deleteCounter int
	userCount     int
)

// UserList aaa
type UserList struct {
	Uid string
	At  string
	Ats string
}

// Handler aaa
func Handler(ctx context.Context) ([]UserList, error) {
	users := getTokens(ctx)

	fmt.Println(users)
	fmt.Println(os.Getenv("CK"))

	anaconda.SetConsumerKey(os.Getenv("CK"))
	anaconda.SetConsumerSecret(os.Getenv("CS"))

	exeCounter = 0
	deleteCounter = 0
	userCount = 0

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(account UserList) {
			defer wg.Done()
			helpTweetDelete(account.At, account.Ats)
		}(user)
	}
	wg.Wait()

	fmt.Println(exeCounter)
	fmt.Println(deleteCounter)
	fmt.Println(userCount)

	return users, nil
}

func main() {
	lambda.Start(Handler)
}

func getTokens(ctx context.Context) []UserList {
	ulget := []UserList{}
	// Dynamoから全権取得
	// if _, err := g.GetAll(datastore.NewQuery("UserList"), &ulget); err != nil {
	// 	fmt.Println(err.Error)
	// }
	return ulget
}

// helpTweetDelete ツイートを取得し救援ツイートを削除
func helpTweetDelete(at string, ats string) (bool, string) {

	api := anaconda.NewTwitterApi(at, ats)
	// api.HttpClient.Transport = trp

	v := url.Values{}
	v.Set("count", "150")

	tweetList, _ := api.GetUserTimeline(v)

	errText := "Error:"
	result := true

	var wgDelete sync.WaitGroup
	for _, t := range tweetList {
		wgDelete.Add(1)
		exeCounter++
		go func(tweet anaconda.Tweet) {
			defer wgDelete.Done()

			gu := `<a href="http://granbluefantasy.jp/" rel="nofollow">グランブルー ファンタジー</a>`
			if tweet.Source == gu {
				_, err := api.DeleteTweet(tweet.Id, false)
				if err != nil {
					errText = errText + err.Error()
					result = false
				}
				deleteCounter++
			}
		}(t)
	}
	wgDelete.Wait()

	return result, errText
}
