cd api-lambda
env GOOS=linux GOARCH=amd64 go build -o /tmp/main
zip -j /tmp/companies-lambda.zip /tmp/main

aws s3 cp /tmp/companies-lambda.zip s3://gofit-lambda-functions/

aws lambda update-function-code --function-name gofit-companies-stack-Lambda-L0NV2RT7FI88 \
--s3-bucket gofit-lambda-functions --s3-key companies-lambda.zip

cd ../custom-authorizer-lambda
zip -r /tmp/api-gateway-authorizer.zip api-gateway-authorizer.py
aws s3 cp /tmp/api-gateway-authorizer.zip s3://gofit-lambda-functions/

aws lambda update-function-code --function-name gofit-companies-stack-Lambda-LY75WETBCIX \
--s3-bucket gofit-lambda-functions --s3-key companies-lambda.zip
