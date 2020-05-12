package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	api *anaconda.TwitterApi
)

func init() {

}

// Response : api gateway you
type Response events.APIGatewayProxyResponse

// Handler : main logic
func Handler(ctx context.Context) (Response, error) {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api = anaconda.NewTwitterApi("", "")

	callbackURL := os.Getenv("CALLBACK_URL")
	fmt.Println(callbackURL)
	url, _, err := api.AuthorizationURL(callbackURL)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// ここからそのまま
	// var buf bytes.Buffer

	// body, err := json.Marshal(map[string]interface{}{
	// 	"message": "Okay so your other function also executed successfully!",
	// 	"tmpCred": tmpCred,
	// 	"url":     url,
	// })
	// if err != nil {
	// 	return Response{StatusCode: 500}, err
	// }
	// json.HTMLEscape(&buf, body)

	html := `<!DOCTYPE html><html>`
	html += `<head>`
	html += `<meta charset="utf-8"/>`
	html += `<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css" rel="stylesheet">`
	html += `<link href='https://fonts.googleapis.com/css?family=Open+Sans:400,700' rel='stylesheet' type='text/css'>`
	html += `<meta name="viewport" content="width=device-width, initial-scale=0.9">`
	html += `<title>グラブルの救援ツイート消すツール</title>`
	html += `</head>`
	html += `<body>`
	html += `<h2> グラブルの救援ツイート消すツール(0.0.1α版) </h2>`
	html += `<br> Twitter認証すると 定期的に ツイートを探索し グラブルの救援ツイート(AP/BP回復ツイート含む) を削除します`
	html += `<br><br> <a href="` + url + `">Twitter ログイン</a>`
	html += `<br>`
	html += `<br> ※運悪くタイミングが被ると救援出した直後に削除されてしまいます。`
	html += `<br> ※登録解除は「設定」「アプリとセッション」から「グラブル救援ツイートクリーナー」を削除してください。`
	html += `</body></html>`

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            html,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
