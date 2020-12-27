package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mtslzr/hlscale/pkg/constants"
	"github.com/mtslzr/hlscale/pkg/exams"
	log "github.com/sirupsen/logrus"
)

// Event is the data passed to the Lambda.
type Event struct {
	Body    string   `json:"body"`
	Records []Record `json:"Records"`
}

// Body is the POST data passed during API calls.
type Body struct {
	Exam     exams.Exam `json:"exam"`
	Function string     `json:"function"`
}

// Record is the data passed from Dynamo stream triggers.
type Record struct {
	Change struct {
		NewImage map[string]*dynamodb.AttributeValue `json:"NewImage"`
	} `json:"dynamodb"`
}

type Handler struct {
	Context context.Context
}

func (h Handler) Handle(ctx context.Context, event Event) (events.APIGatewayProxyResponse, error) {
	if len(event.Records) > 0 {
		fmt.Printf("change -> %+v\n", event.Records[0].Change)
		log.Infof("Running %s...", constants.CreateEvent)
		return sendResponse(nil)
	} else {
		var body Body
		err := json.Unmarshal([]byte(event.Body), &body)
		if err != nil {
			return sendResponse(err)
		}

		switch body.Function {
		case constants.CreateExam:
			log.Infof("Running %s...", constants.CreateExam)
			return sendResponse(exams.CreateExam(body.Exam))
		case constants.StartScale:
			log.Infof("Running %s...", constants.StartScale)
		case constants.EndScale:
			log.Infof("Running %s...", constants.EndScale)
		}
		return sendResponse(errors.New("unknown or missing function"))
	}
}

func sendResponse(err error) (events.APIGatewayProxyResponse, error) {
	var resp events.APIGatewayProxyResponse
	if err == nil {
		resp = events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            "Success!",
			IsBase64Encoded: false,
		}
	} else {
		resp = events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "Error!",
			IsBase64Encoded: false,
		}
	}
	return resp, err
}
