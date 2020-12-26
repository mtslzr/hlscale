package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mtslzr/hlscale/pkg/constants"
	"github.com/mtslzr/hlscale/pkg/exams"
	log "github.com/sirupsen/logrus"
)

type Body struct {
	Exam     exams.Exam `json:"exam"`
	Function string     `json:"function"`
}

type Handler struct {
	Context context.Context
}

func (h Handler) Handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body Body
	err := json.Unmarshal([]byte(event.Body), &body)
	if err != nil {
		return sendResponse(err)
	}

	switch body.Function {
	case constants.CreateExam:
		log.Infof("Running %s...", constants.CreateExam)
		return sendResponse(exams.CreateExam(body.Exam))
	case constants.CreateEvent:
		log.Infof("Running %s...", constants.CreateEvent)
	case constants.StartScale:
		log.Infof("Running %s...", constants.StartScale)
	case constants.EndScale:
		log.Infof("Running %s...", constants.EndScale)
	}
	return sendResponse(errors.New("unknown or missing function"))
}

func sendResponse(err error) (events.APIGatewayProxyResponse, error) {
	var resp events.APIGatewayProxyResponse
	if err == nil {
		resp = events.APIGatewayProxyResponse{
			StatusCode:        200,
			Body:              "Success!",
			IsBase64Encoded:   false,
		}
	} else {
		resp = events.APIGatewayProxyResponse{
			StatusCode:        500,
			Body:              "Error!",
			IsBase64Encoded:   false,
		}
	}
	return resp, err
}