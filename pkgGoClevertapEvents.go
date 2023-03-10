package pkgGoClevertapEvents

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// "github.com/aws/aws-sdk-go-v2/aws"

type ClevertapEventSender struct {
	config   *aws.Config
	client   eventbridge.Client
	eventBus string
	source   string
}

func (c *ClevertapEventSender) Initialize(config *aws.Config, eventBus string, source string) {
	c.config = config
	c.client = *eventbridge.NewFromConfig(*c.config)
	c.eventBus = eventBus
	c.source = source
}

func (c *ClevertapEventSender) SendEventToEc(ctx context.Context, data ClevertapEventPayload, detail_type string) {
	var inputList []types.PutEventsRequestEntry
	detail, _ := json.Marshal(data)
	inputList = append(inputList, defineEventHash(string(detail), detail_type, c.source, c.eventBus))
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
		Detail:       aws.String(detail),
		DetailType:   aws.String(detail_type),
		EventBusName: aws.String(eventBus),
		Resources:    []string{},
		Source:       &source,
		Time:         &curtime,
		TraceHeader:  &curTimeStr,
	}
}

type ClevertapEventData struct {
	Type        string                 `json:"type"`
	EventName   string                 `json:"evtName"`
	Identity    string                 `json:"identity" binding:"required"`
	Time        int64                  `json:"ts"`
	EventData   map[string]interface{} `json:"evtData"`
	ProfileData map[string]interface{} `json:"profileData"`
}

type ClevertapEventPayload struct {
	D []ClevertapEventData `json:"d"`
}
