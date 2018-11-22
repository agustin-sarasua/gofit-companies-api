First, remember to create the bucket "gofit-lambda-functions" to store the lambda functions.

Package and Deploy the lambdas

1. api-lambda 
2. Custom Authorizer Lambda


aws dynamodb query \
    --table-name Companies \
    --key-condition-expression "UserSub = :name" \
    --expression-attribute-values  '{":name":{"S":"776d21e0-3b27-49df-a878-e0c7458c3100"}}'