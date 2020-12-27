package constants

const (
	// Events
	EventsArn       = "arn:aws:iam::645714156459:role/hlscale-cwevents"
	EventsLambdaArn = "arn:aws:lambda:us-east-1:645714156459:function:hlscale"

	// Exams
	ExamsTableName = "hlscale-exams"

	// Handler
	CreateEvent = "createEvents"
	CreateExam  = "createExam"
	EndScale    = "endScale"
	StartScale  = "startScale"
)
