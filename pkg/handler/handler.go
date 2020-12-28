package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mtslzr/hlscale/pkg/cwevents"
	"github.com/mtslzr/hlscale/pkg/scale"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mtslzr/hlscale/pkg/constants"
	"github.com/mtslzr/hlscale/pkg/exams"
	log "github.com/sirupsen/logrus"
)

// Body is the POST data passed during API calls.
type Body struct {
	Event    cwevents.Event `json:"event"`
	Exam     exams.Exam     `json:"exam"`
	Function string         `json:"function"`
}

// Event is the data passed to the Lambda.
type Event struct {
	Body    string `json:"body"`
	Records []struct {
		Change struct {
			NewImage map[string]*dynamodb.AttributeValue `json:"NewImage"`
		} `json:"dynamodb"`
	} `json:"Records"`
}

// Handler is the context for executing Lambda functions.
type Handler struct {
	Context context.Context
}

// Handle is a pseudo-router to set up the function to run.
func (h Handler) Handle(ctx context.Context, event Event) (events.APIGatewayProxyResponse, error) {
	if len(event.Records) > 0 {
		log.Infof("Running %s...", constants.CreateEvent)
		record, err := cwevents.ParseRecord(event.Records[0].Change.NewImage)
		if err != nil {
			return sendResponse(err)
		}
		return sendResponse(cwevents.CreateEvents(record))
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
			if err = scale.UpdateCapacity(body.Event.Students); err != nil {
				return sendResponse(err)
			}
			return sendResponse(cwevents.DeleteRule(body.Event.Name))
		case constants.EndScale:
			log.Infof("Running %s...", constants.EndScale)
			if err = scale.UpdateCapacity(0); err != nil {
				return sendResponse(err)
			}
			return sendResponse(cwevents.DeleteRule(body.Event.Name))

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
