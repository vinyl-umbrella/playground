package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

func HandleRequest(ctx context.Context) {
	env := os.Getenv("ENV")
	fmt.Println("ENV: ", env)
}
