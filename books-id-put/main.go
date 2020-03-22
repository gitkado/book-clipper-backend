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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// JSONリクエストからstructオブジェクト作成
	r := new(Request)
	err := json.Unmarshal([]byte(request.Body), r)
	if err != nil {
		panic(err.Error())
	}

	// 登録用structオブジェクト準備
	time_now := time.Now().UTC()
	r.Book.UpdatedAt = time_now

	// jsonレスポンス用に変換
	b_byte, _ := json.Marshal(r.Book)

	// dynamo検索用
	created_at := r.Book.CreatedAt

	// dynamoにupdate
	u := table.Update("created_at", created_at)
	u.Set("title", r.Book.Title)
	u.Set("url", r.Book.Url)
	u.SetSet("tag", r.Book.Tag)
	u.Set("is_book", r.Book.IsBook)
	u.Set("is_ebook", r.Book.IsEbook)
	u.Set("updated_at", r.Book.UpdatedAt)

	// dynamoにput
	err = u.Run()
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
