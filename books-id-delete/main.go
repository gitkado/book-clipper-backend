package main

// https://github.com/guregu/dynamo
import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// RequestBody用Struct定義
type Request struct {
	CreatedAt string `dynamo:"created_at" json:"created_at"`
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
	// dynamoからdelete
	created_at := request.CreatedAt
	if err := table.Delete("created_at", created_at).Run(); err != nil {
		panic(err.Error())
	}

	// jsonレスポンス用に変換
	b_byte, _ := json.Marshal(request)

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
