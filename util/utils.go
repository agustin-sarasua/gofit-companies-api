package util

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func GetClaimsSub(ctx events.APIGatewayProxyRequestContext) string {
	jc, _ := json.Marshal(ctx.Authorizer)
	fmt.Print(string(jc))
	r := make(map[string]interface{})
	err := json.Unmarshal(jc, &r)
	if err != nil {
		fmt.Printf("Something went wrong %v", err)
	}
	return r["principalId"].(string)
}
