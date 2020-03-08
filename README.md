# book-clipper-backend

## Language
- go

## Deploy
- Amazon API Gaetway
- Amazon DynamoDB
- AWS Lambda
- CodeBuild(TODO)

## Run

```bash
GOOS=linux GOARCH=amd64 go build -o books-post/books-post ./books-post
GOOS=linux GOARCH=amd64 go build -o books-get/books-get ./books-get
GOOS=linux GOARCH=amd64 go build -o books-id-get/books-id-get ./books-id-get
GOOS=linux GOARCH=amd64 go build -o books-id-put/books-id-put ./books-id-put
GOOS=linux GOARCH=amd64 go build -o books-id-delete/books-id-delete ./books-id-delete

sam local invoke BooksPostFunction --event events/books-post.json
sam local invoke BooksGetFunction --no-event
sam local invoke BooksIdGetFunction --event events/books-id-get.json
sam local invoke BooksIdPutFunction --event events/books-id-put.json
sam local invoke BooksIdDeleteFunction --event events/books-id-delete.json
```

## URL
- https://${ServerlessRestApi}.execute-api.${AWS::Region}.amazonaws.com/Prod/books
