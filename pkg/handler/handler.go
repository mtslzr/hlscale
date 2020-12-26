package handler

import (
	"context"
	"github.com/mtslzr/hlscale/pkg/constants"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	Function string `json:"function"`
}

type Handler struct {
	Context context.Context
}

func (h Handler) Handle(ctx context.Context, event Event) error {
	switch event.Function{
	case constants.CreateExam:
		log.Infof("Running %s...", constants.CreateExam)
	case constants.CreateEvent:
		log.Infof("Running %s...", constants.CreateEvent)
	case constants.StartScale:
		log.Infof("Running %s...", constants.StartScale)
	case constants.EndScale:
		log.Infof("Running %s...", constants.EndScale)
	}
	return nil
}