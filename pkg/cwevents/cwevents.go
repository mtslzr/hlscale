package cwevents

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mtslzr/hlscale/pkg/constants"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	RuleArn  string `json:"rule_arn"`
	Students int    `json:"students"`
}

type EventInput struct {
	Event    Event  `json:"event"`
	Function string `json:"function"`
}

type Record struct {
	Name     string `json:"name"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Students int    `json:"students"`
}

// ParseRecord converts a Dynamo stream record to a Record object.
func ParseRecord(record map[string]*dynamodb.AttributeValue) (Record, error) {
	students, err := strconv.Atoi(*record["students"].N)
	if err != nil {
		return Record{}, err
	}

	return Record{
		Name:     *record["name"].S,
		Start:    *record["start"].N,
		End:      *record["end"].N,
		Students: students,
	}, nil
}

// CreateEvents adds CloudWatch events at start and end of exam.
func CreateEvents(record Record) error {
	sess, err := session.NewSession()
	if err != nil {
		log.Errorf("Error creating AWS session: %s", err)
		return err
	}

	svc := cloudwatchevents.New(sess)
	shortName := strings.ReplaceAll(record.Name, " ", "-")

	startName := fmt.Sprintf("Scale-Up-%s", shortName)
	startRule := putRule(startName, record.Start)
	startResp, err := svc.PutRule(&startRule)
	if err != nil {
		log.Errorf("Error creating rule: %s", err)
		return err
	}

	startInput, err := json.Marshal(EventInput{
		Event:    Event{
			RuleArn:  *startResp.RuleArn,
			Students: record.Students,
		},
		Function: constants.StartScale,
	})
	if err != nil {
		log.Errorf("Error creating event input: %s", err)
		return err
	}

	startTarget := putTarget(startName, startInput)
	_, err = svc.PutTargets(&startTarget)
	if err != nil {
		log.Errorf("Error creating target: %s", err)
		return err
	}

	endName := fmt.Sprintf("Scale-Down-%s", shortName)
	endRule := putRule(endName, record.End)
	endResp, err := svc.PutRule(&endRule)
	if err != nil {
		log.Errorf("Error creating rule: %s", err)
		return err
	}

	endInput, err := json.Marshal(EventInput{
		Event:    Event{
			RuleArn:  *endResp.RuleArn,
		},
		Function: constants.EndScale,
	})
	if err != nil {
		log.Errorf("Error creating event input: %s", err)
		return err
	}

	endTarget := putTarget(endName, endInput)
	_, err = svc.PutTargets(&endTarget)
	if err != nil {
		log.Errorf("Error creating target: %s", err)
		return err
	}

	return nil
}

func putRule(name string, sched string) cloudwatchevents.PutRuleInput {
	return cloudwatchevents.PutRuleInput{
		Name:               aws.String(name),
		RoleArn:            aws.String(constants.EventsArn),
		ScheduleExpression: aws.String(unixToCron(sched)),
	}
}

func putTarget(name string, input []byte) cloudwatchevents.PutTargetsInput {
	return cloudwatchevents.PutTargetsInput{
		Rule: aws.String(name),
		Targets: []*cloudwatchevents.Target{
			{
				Arn:   aws.String(constants.EventsLambdaArn),
				Id:    aws.String(strings.ReplaceAll(name, "-", "")),
				Input: aws.String(string(input)),
			},
		},
	}
}

func unixToCron(ts string) string {
	tsInt, _ := strconv.Atoi(ts)
	input := time.Unix(int64(tsInt), 0)
	cron := []string{
		strconv.Itoa(input.Minute()),
		strconv.Itoa(input.Hour()),
		strconv.Itoa(input.Day()),
		strconv.Itoa(int(input.Month())),
		"?",
		strconv.Itoa(input.Year()),
	}
	return fmt.Sprintf("cron(%s)", strings.Join(cron, " "))
}
