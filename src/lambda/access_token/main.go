package main

import (
	"context"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	api *anaconda.TwitterApi
)

func init() {

}

// Response : api gateway you
type Response events.APIGatewayProxyResponse

// Handler : main logic
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api = anaconda.NewTwitterApi("", "")

	oauthToken := request.QueryStringParameters["oauth_token"]
	oauthVerifier := request.QueryStringParameters["oauth_verifier"]

	// credentials 疑似生成
	api.Credentials.Token = oauthToken

	c, _, err := api.GetCredentials(api.Credentials, oauthVerifier)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// ユーザid取得
	api := anaconda.NewTwitterApi(c.Token, c.Secret)
	me, _ := api.GetSelf(url.Values{})

	// DBへユーザ情報保存
	err = putTokens(me.IdStr, c.Token, c.Secret)
	if err != nil {
		return nil, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            me.IdStr + "|||" + c.Token + "|||" + c.Secret,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func putTokens(id string, token string, tokenSeaclet string) error {
	userData := UserData{
		UID:    id,
		Token:  token,
		Secret: tokenSeaclet,
	}
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	table := db.Table(os.Getenv("DB_NAME"))
	err := table.Put(userData).Run()
	if err != nil {
		return err
	}
	return nil
}

type UserData struct {
	UID    string `dynamo:"UserID"`
	Token  string `dynamo:"Token"`
	Secret string `dynamo:"Secret"`
}
