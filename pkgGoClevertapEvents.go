package pkgGoClevertapEvents

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// "github.com/aws/aws-sdk-go-v2/aws"

type ClevertapEventSender struct {
	config   *aws.Config
	client   eventbridge.Client
	eventBus string
}

func (c *ClevertapEventSender) Initialize(config *aws.Config, eventBus string) {
	c.config = config
	c.client = *eventbridge.NewFromConfig(*c.config)
	c.eventBus = eventBus
}

func (c *ClevertapEventSender) SendEventToEc(ctx context.Context, detail string, detail_type string, source string) {
	var inputList []types.PutEventsRequestEntry
	inputList = append(inputList, defineEventHash(detail, detail_type, source, c.eventBus))
	input := eventbridge.PutEventsInput{
		Entries: inputList,
	}
	out, err := c.client.PutEvents(ctx, &input)
	fmt.Println(out, err)
}

func defineEventHash(detail string, detail_type string, source string, eventBus string) types.PutEventsRequestEntry {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	curtime := (time.Now().In(loc))
	curTimeStr := curtime.String() + source
	return types.PutEventsRequestEntry{
		Detail:       &detail,
		DetailType:   &detail_type,
		EventBusName: &eventBus,
		Resources:    []string{},
		Source:       &source,
		Time:         &curtime,
		TraceHeader:  &curTimeStr,
	}
}

type ClevertapEventData struct {
	Type        string            `json:"type"`
	EventName   string            `json:"evtName"`
	Identity    string            `json:"identity" binding:"required"`
	Time        int64             `json:"ts"`
	EventData   map[string]string `json:"evtData"`
	ProfileData map[string]string `json:"profileData"`
}

type ClevertapEventPayload struct {
	D []ClevertapEventData `json:"d"`
}
