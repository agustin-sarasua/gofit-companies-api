package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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

func createCompany(c *gin.Context) {
	apiGwContext, _ := ginLambda.GetAPIGatewayContext(c.Request)
	e := Company{Rol: RolOwner}

	err := c.BindJSON(&e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	uid, _ := uuid.NewV4()
	e.ID = uid.String()

	startTime := time.Now().Format("2006-01-02T15:04:05")
	e.Timestamp = startTime
	if e.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Name not found",
		})
		return
	}
	e.UserSub = getClaimsSub(apiGwContext)
	e.Status = StatusActive
	err = putCompany(&e)
	if err != nil {
		fmt.Printf("Error saving item in db %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, e)
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
	err = putStaff(&e)
	if err != nil {
		fmt.Printf("Error saving item in db %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, e)
}

func listCompanies(c *gin.Context) {
	apiGwContext, _ := ginLambda.GetAPIGatewayContext(c.Request)
	userSub := getClaimsSub(apiGwContext)
	cs, err := getUserCompanies(userSub, 10)
	if err != nil {
		fmt.Printf("Error saving item in db %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	log.Printf("Length of cs %d", len(cs))

	c.JSON(http.StatusOK, &struct {
		Companies []*Company `json:"Companies"`
	}{Companies: cs})
}

func getCompanyData(c *gin.Context) {
	companyID := c.Param("id")
	co, err := getCompanyWithStaff(companyID, 100)
	if err != nil {
		fmt.Printf("Error saving item in db %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, co)
}

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()
		r.GET("/companies", listCompanies)
		r.GET("/companies/:id", getCompanyData)
		r.POST("/companies", createCompany)
		r.POST("/companies/:id/staff", createStaff)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
