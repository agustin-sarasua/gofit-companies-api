cd api-lambda
env GOOS=linux GOARCH=amd64 go build -o /tmp/main
zip -j /tmp/companies-lambda.zip /tmp/main
aws s3 cp /tmp/companies-lambda.zip s3://gofit-lambda-functions/

cd ../custom-authorizer-lambda
zip -j /tmp/api-gateway-authorizer.zip api-gateway-authorizer.py
aws s3 cp /tmp/api-gateway-authorizer.zip s3://gofit-lambda-functions/