package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var infoLogger = log.New(os.Stdout, "INFO ", log.Llongfile)

func getClaimsSub(ctx events.APIGatewayProxyRequestContext) string {
	jc, _ := json.Marshal(ctx.Authorizer)
	fmt.Print(string(jc))
	r := make(map[string]interface{})
	err := json.Unmarshal(jc, &r)
	if err != nil {
		fmt.Printf("Something went wrong %v", err)
	}
	return r["principalId"].(string)
}

// TODO validate if the company exists and if the UserSub is the owner of it.
// Send push notification to Staff to accept being Staff
// Once accepted he is added
func createStaff(c *gin.Context) {
	companyID := c.Param("id")
	e := Staff{}
	err := c.BindJSON(&e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	if e.Rol == "" {
		e.Rol = RolStaff
	}
	apiGwContext, _ := ginLambda.GetAPIGatewayContext(c.Request)
	// TODO validate UserSub exists
	e.CreatedBy = getClaimsSub(apiGwContext)
	e.Status = StatusPending
	e.CompanyID = companyID
	e.SortKey = fmt.Sprintf("staff-%s", e.UserSub)
	err = putStaff(&e)
	if err != nil {
		fmt.Printf("Error saving item in db %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, e)
}

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()
		r.POST("/companies/:id/staff", createStaff)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
