package main

import (
	"bytes"
	"context"
	"io"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	data := map[string]interface{}{
		"Name": "...",
	}
	render("./views/home.html", w, data)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func render(v string, w io.Writer, data map[string]interface{}) {
	tmpl := template.Must(template.ParseFiles("./views/layout.html", v))
	tmpl.Execute(w, data)
}
