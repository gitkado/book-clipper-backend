# book-clipper-backend

## Language
- go

## Deploy
- Amazon API Gaetway
- Amazon DynamoDB
- AWS Lambda
- AWS CodeBuild

## Run

```sh
GOOS=linux GOARCH=amd64 go build -o books-post/books-post ./books-post
GOOS=linux GOARCH=amd64 go build -o books-get/books-get ./books-get
GOOS=linux GOARCH=amd64 go build -o books-id-get/books-id-get ./books-id-get
GOOS=linux GOARCH=amd64 go build -o books-id-put/books-id-put ./books-id-put
GOOS=linux GOARCH=amd64 go build -o books-id-delete/books-id-delete ./books-id-delete

sam local start-api
curl 127.0.0.1:3000/books
curl 127.0.0.1:3000/books/{created_at}
curl -X POST -H "Content-Type: application/json" -d '{"book": {"title": "Nuxtjs Tutorial","url": "","tag": ["Nuxt"],"is_book": true,"is_ebook": false}}' 127.0.0.1:3000/books
curl -X PUT -H "Content-Type: application/json" -d '{"book": {"title": "Vue.js&Nuxt.js Tutorial","url": "","tag": ["Vue","Nuxt"],"is_book": true,"is_ebook": true,"created_at": "{created_at}"}}' 127.0.0.1:3000/books/{created_at}
curl -X DELETE 127.0.0.1:3000/books/{created_at}
```

## Build

```sh
sam package \
    --template-file template.yaml \
    --s3-bucket {S3BucketName} \
    --output-template-file packaged-template.yaml \
    --region ap-northeast-1
sam deploy \
    --template-file packaged-template.yaml \
    --stack-name cfn-lambda-bookclipper \
    --capabilities CAPABILITY_IAM
```

## URL
- https://${ServerlessRestApi}.execute-api.${AWS::Region}.amazonaws.com/Prod/books
