package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	api           *anaconda.TwitterApi
	exeCounter    int
	deleteCounter int
	userCount     int
)

func init() {

}

// Handler : main logic
func Handler(ctx context.Context) error {

	// 初期設定
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api = anaconda.NewTwitterApi("", "")

	// 削除対象ユーザデータを取得
	users, err := getTokens()
	if err != nil {
		return err
	}

	// 削除ループ用カウンタ
	exeCounter = 0
	deleteCounter = 0
	userCount = 0

	// 地獄のループ
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(account UserData) {
			defer wg.Done()
			helpTweetDelete(account.Token, account.Secret)
			fmt.Println(account.UID + "Done")
		}(user)
	}
	wg.Wait()
	fmt.Println(strconv.Itoa(exeCounter))
	fmt.Println(strconv.Itoa(deleteCounter))

	return nil
}

func main() {
	lambda.Start(Handler)
}

type UserData struct {
	UID    string `dynamo:"UserID"`
	Token  string `dynamo:"Token"`
	Secret string `dynamo:"Secret"`
}

// helpTweetDelete ツイートを取得し救援ツイートを削除
func helpTweetDelete(at string, ats string) (bool, string) {

	api := anaconda.NewTwitterApi(at, ats)

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

func getTokens() ([]UserData, error) {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	table := db.Table(os.Getenv("DB_NAME"))
	var results []UserData
	err := table.Scan().All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
