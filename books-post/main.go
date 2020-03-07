package main

// https://github.com/google/uuid
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

type RequestBook struct {
	Title   string   `json:"title"`
	Url     string   `json:"url"`
	Tag     []string `json:"tag"`
	IsBook  bool     `json:"is_book"`
	IsEbook bool     `json:"is_ebook"`
}

// RequestBody用Struct定義
type Request struct {
	Book RequestBook `json:"book"`
}

// DynamoDB/Book用Struct定義
type Book struct {
	Title     string    `dynamo:"title" json:"title"`
	Url       string    `dynamo:"url" json:"url"`
	Tag       []string  `dynamo:"tag,set" json:"tag"`
	IsBook    bool      `dynamo:"is_book" json:"is_book"`
	IsEbook   bool      `dynamo:"is_ebook" json:"is_ebook"`
	CreatedAt time.Time `dynamo:"created_at" json:"created_at"`
	UpdatedAt time.Time `dynamo:"updated_at" json:"updated_at"`
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
	time_now := time.Now().UTC()
	b := Book{
		Title:     request.Book.Title,
		Url:       request.Book.Url,
		Tag:       request.Book.Tag,
		IsBook:    request.Book.IsBook,
		IsEbook:   request.Book.IsEbook,
		CreatedAt: time_now,
		UpdatedAt: time_now,
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
