package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mtslzr/hlscale/pkg/handler"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	h := handler.Handler{}
	lambda.Start(h.Handle)
}