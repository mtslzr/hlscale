package scale

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/mtslzr/hlscale/pkg/constants"
	log "github.com/sirupsen/logrus"
	"math"
)

// UpdateCapacity turns ASG capacity up to pre-warm an exam.
func UpdateCapacity(num int) error {
	sess, err := session.NewSession()
	if err != nil {
		log.Errorf("Error creating AWS session: %s", err)
		return err
	}

	svc := autoscaling.New(sess)
	_, err = svc.SetDesiredCapacity(&autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(constants.ScaleGroupName),
		DesiredCapacity:      aws.Int64(getCapacity(num)),
	})
	if err != nil {
		log.Errorf("Error setting ASG capacity: %s", err)
	}
	return err
}

func getCapacity(students int) int64 {
	if capacity := int64(math.Ceil(float64(students / 100))); capacity < 1 {
		return 1
	} else {
		return capacity
	}
}