package main

// https://github.com/guregu/dynamo
import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// RequestBody用Struct定義
type Request struct {
	Title   string   `json:"title"`
	Url     string   `json:"url"`
	Tag     []string `json:"tag"`
	IsBook  bool     `json:"is_book"`
	IsEbook bool     `json:"is_ebook"`
}

// DynamoDB/Book用Struct定義
type Book struct {
	Title     string    `dynamo:"title" json:"title"`
	Url       string    `dynamo:"url" json:"url"`
	Tag       []string  `dynamo:"tag,set" json:"tag"`
	IsBook    bool      `dynamo:"is_book" json:"is_book"`
	IsEbook   bool      `dynamo:"is_ebook" json:"is_ebook"`
	Timestamp time.Time `dynamo:"timestamp" json:"timestamp"`
}

// 変数定義
var (
	// dynamoオブジェクト作成
	db = dynamo.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	// tableオブジェクト作成
	table = db.Table("book-clipper")
)

func Handler(request Request) (events.APIGatewayProxyResponse, error) {
	// 登録用structオブジェクト作成
	b := Book{
		Title:     request.Title,
		Url:       request.Url,
		Tag:       request.Tag,
		IsBook:    request.IsBook,
		IsEbook:   request.IsEbook,
		Timestamp: time.Now().UTC(),
	}
	// jsonレスポンス用に変換
	b_byte, _ := json.Marshal(b)

	// dynamoにput
	if err := table.Put(b).Run(); err != nil {
		panic(err.Error())
	}

	// lambdaレスポンス返却
	return events.APIGatewayProxyResponse{
		Body:            string(b_byte),
		StatusCode:      200,
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
