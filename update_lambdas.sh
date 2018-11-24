cd api-lambda
env GOOS=linux GOARCH=amd64 go build -o /tmp/main
zip -j /tmp/companies-lambda.zip /tmp/main

aws lambda update-function-code --function-name GoFitCompaniesApi \
--zip-file fileb:///tmp/companies-lambda.zip


zip -r /tmp/api-gateway-authorizer.zip  *
