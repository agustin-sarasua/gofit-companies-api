package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/agustin-sarasua/gofit-companies-api/model"
	"github.com/agustin-sarasua/gofit-companies-api/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var ginLambda *ginadapter.GinLambda

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var infoLogger = log.New(os.Stdout, "INFO ", log.Llongfile)

func createCompany(c *gin.Context) {
	apiGwContext, _ := ginLambda.GetAPIGatewayContext(c.Request)
	e := model.Company{Rol: model.RolOwner}

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
	e.UserSub = util.GetClaimsSub(apiGwContext)
	e.Status = model.StatusActive
	e.SortKey = fmt.Sprintf("company-%s", e.ID)
	e.Staff = nil

	err = putCompany(&e)
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
	userSub := util.GetClaimsSub(apiGwContext)
	cs, err := getUserCompanies(userSub, 10)
	if err != nil {
		fmt.Printf("Error saving item in db %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	log.Printf("Length of cs %d", len(cs))

	c.JSON(http.StatusOK, &struct {
		Companies []*model.Company `json:"Companies"`
	}{Companies: cs})
}

func deleteCompany(c *gin.Context) {
	companyID := c.Param("id")
	co, err := getCompanyWithStaff(companyID, 100)
	if err != nil {
		fmt.Printf("Error deleting company %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, co)
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
		r.DELETE("/companies/:id", deleteCompany)
		r.GET("/companies", listCompanies)
		r.GET("/companies/:id", getCompanyData)
		r.POST("/companies", createCompany)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
