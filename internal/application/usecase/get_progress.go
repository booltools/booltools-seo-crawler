package usecase

import (
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/queue"
)

type GetProgressUseCase struct {
	progressBroker *queue.ProgressBroker
}

func NewGetProgressUseCase(progressBroker *queue.ProgressBroker) *GetProgressUseCase {
	return &GetProgressUseCase{
		progressBroker: progressBroker,
	}
}

func (uc *GetProgressUseCase) Subscribe(jobID string) chan queue.ProgressEvent {
	return uc.progressBroker.Subscribe(jobID)
}

func (uc *GetProgressUseCase) Unsubscribe(jobID string, channel chan queue.ProgressEvent) {
	uc.progressBroker.Unsubscribe(jobID, channel)
}
