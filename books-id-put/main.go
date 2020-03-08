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

type RequestOld struct {
	CreatedAt string `json:"created_at"`
}

// RequestBody用Struct定義
type Request struct {
	Book Book       `json:"book"`
	Old  RequestOld `json:"old"`
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
	// jsonレスポンス用に変換
	b_byte, _ := json.Marshal(request.Book)

	// dynamoからget
	created_at := request.Old.CreatedAt

	// dynamoにupdate
	u := table.Update("created_at", created_at)
	u.Set("title", request.Book.Title)
	u.Set("url", request.Book.Url)
	u.SetSet("tag", request.Book.Tag)
	u.Set("is_book", request.Book.IsBook)
	u.Set("is_ebook", request.Book.IsEbook)
	u.Set("updated_at", time.Now().UTC())

	// dynamoにput
	err := u.Run()
	if err != nil {
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
