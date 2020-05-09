package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	url, tmpCred, err := api.AuthorizationURL(callbackURL)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// ここからそのまま
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Okay so your other function also executed successfully!",
		"tmpCred": tmpCred,
		"url":     url,
	})
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	json.HTMLEscape(&buf, body)

	html := `<!DOCTYPE html><html><head><meta charset="utf-8"/></head><body>`
	html += `<br> <a href="` + url + `">Twitter ログイン</a> <br>`
	// html += buf.String()
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
