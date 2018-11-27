First, remember to create the bucket "gofit-lambda-functions" to store the lambda functions.

Package and Deploy the lambdas

1. api-lambda 
2. Custom Authorizer Lambda


DynamoDB Table

- Composite key: Partition Key + Range Key

CompanyID   |   SortKey   | Attributes
[CompanyID] | Service-[ServiceID]#Class-[ClassID] | [UserSub]
[CompanyID] | Staff-[StaffID] | [UserSub]

Table
company1 | company-company1 | user1 | owner
company1 | staff-user2 | user2 | staff
company1 | staff-user3 | user3 | staff
company2 | company-company2 | user2 | owner
company2 | staff-user1 | user1 | staff

Global Secondary index
user1 | company-company1 
user2 | staff-user2 
user3 | staff-user3 
user2 | company-company2 
user1 | staff-user1 



Global Seconday Index
[UserSub] | [CompanyID]

GET /companies
Query GSI by UserSub

GET /companies/:id



aws dynamodb query \
    --table-name Companies \
    --key-condition-expression "UserSub = :name" \
    --expression-attribute-values  '{":name":{"S":"776d21e0-3b27-49df-a878-e0c7458c3100"}}'