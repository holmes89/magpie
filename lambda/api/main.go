package main

import (
	"context"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/holmes89/magpie/lib/database"
	v1 "github.com/holmes89/magpie/lib/handlers/rest/v1"
	"github.com/holmes89/magpie/lib/resource"
)

var (
	muxAdapter *gorillamux.GorillaMuxAdapter
	once       sync.Once
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxAdapter.ProxyWithContext(ctx, request)
}

func setup() {
	router := mux.NewRouter()
	svc := resource.NewService(database.NewConnection())
	v1.MakeV1SiteHandler(router, svc)
	muxAdapter = gorillamux.New(router)
}

func main() {
	once.Do(setup)
	lambda.Start(handleRequest)
}
