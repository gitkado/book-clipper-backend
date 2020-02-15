package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#UpdateItemInput
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#example_DynamoDB_UpdateItem_shared00

func hoge() (events.APIGatewayProxyResponse, error) {
	svc := dynamodb.New(session.New())
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#AT": aws.String("AlbumTitle"),
			"#Y":  aws.String("Year"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String("Louder Than Ever"),
			},
			":y": {
				N: aws.String("2015"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("aaa"),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String("dynamodbtest"),
		UpdateExpression: aws.String("SET #Y = :y, #AT = :t"),
	}

	svc.UpdateItem(input)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string("aaa")),
		StatusCode: 200,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string(ip)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(hoge)
	// lambda.Start(handler)
}
