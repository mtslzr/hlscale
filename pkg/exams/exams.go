package exams

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/mtslzr/hlscale/pkg/constants"
	log "github.com/sirupsen/logrus"
)

type Exam struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Students int    `json:"students"`
}

//CreateExam adds a new Exam to DynamoDB.
func CreateExam(exam Exam) error {
	sess, err := session.NewSession()
	if err != nil {
		log.Errorf("Error creating AWS session: %s", err)
		return err
	}

	exam.ID = uuid.New().String()
	av, err := dynamodbattribute.MarshalMap(exam)
	if err != nil {
		log.Errorf("Error formatting Dynamo item: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(constants.ExamsTableName),
		Item:      av,
	}

	svc := dynamodb.New(sess)
	_, err = svc.PutItem(input)
	if err != nil {
		log.Errorf("Error putting item in Dynamo: %s", err)
	}
	return err
}
