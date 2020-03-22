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

// RequestBody用Struct定義
type Request struct {
	Book Book `json:"book"`
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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// JSONリクエストからstructオブジェクト作成
	r := new(Request)
	err := json.Unmarshal([]byte(request.Body), r)
	if err != nil {
		panic(err.Error())
	}

	// 登録用structオブジェクト作成
	time_now := time.Now().UTC()
	b := Book{
		Title:     r.Book.Title,
		Url:       r.Book.Url,
		Tag:       r.Book.Tag,
		IsBook:    r.Book.IsBook,
		IsEbook:   r.Book.IsEbook,
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
