package queue

import (
	"sync"
)

type PageIssue struct {
	RuleKey  string `json:"ruleKey"`
	Severity string `json:"severity"`
	Result   string `json:"result"`
	Message  string `json:"message"`
	Details  string `json:"details,omitempty"`
}

type ProgressEvent struct {
	JobID           string      `json:"jobId"`
	Status          string      `json:"status"`
	PagesCrawled    int         `json:"pagesCrawled"`
	TotalDiscovered int         `json:"totalDiscovered"`
	IssuesFound     int         `json:"issuesFound"`
	CurrentURL      string      `json:"currentUrl,omitempty"`
	Message         string      `json:"message,omitempty"`
	PageIssues      []PageIssue `json:"pageIssues,omitempty"`
}

type ProgressBroker struct {
	subscribers map[string][]chan ProgressEvent
	mutex       sync.RWMutex
}

func NewProgressBroker() *ProgressBroker {
	return &ProgressBroker{
		subscribers: make(map[string][]chan ProgressEvent),
	}
}

func (b *ProgressBroker) Subscribe(jobID string) chan ProgressEvent {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	channel := make(chan ProgressEvent, 100)
	b.subscribers[jobID] = append(b.subscribers[jobID], channel)
	return channel
}

func (b *ProgressBroker) Unsubscribe(jobID string, channel chan ProgressEvent) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	channels := b.subscribers[jobID]
	for i, subscribedChannel := range channels {
		if subscribedChannel == channel {
			b.subscribers[jobID] = append(channels[:i], channels[i+1:]...)
			close(channel)
			break
		}
	}

	if len(b.subscribers[jobID]) == 0 {
		delete(b.subscribers, jobID)
	}
}

func (b *ProgressBroker) Publish(event ProgressEvent) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	channels := b.subscribers[event.JobID]
	for _, channel := range channels {
		select {
		case channel <- event:
		default:
		}
	}
}
