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
// 検索やソート機能を実装する場合に使用
// type Request struct {
// }

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

func Handler() (events.APIGatewayProxyResponse, error) {
	// 登録用structオブジェクト作成
	var results []Book

	// dynamoからget
	if err := table.Scan().All(&results); err != nil {
		panic(err.Error())
	}

	// jsonレスポンス用に変換
	b_byte, _ := json.Marshal(results)

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
