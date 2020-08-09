package lib

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	clogs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/ayoul3/trailSniffer/cloudwatchlogs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func LookupLogs(client *cloudwatchlogs.Client, logGroup, filter string, startDate, endDate time.Time) (result []*clogs.FilteredLogEvent, err error) {
	log.Infof("Looking up events between %s and %s", startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))

	if result, err = client.FetchLogs(logGroup, filter, startDate, endDate); err != nil {
		if len(result) == 0 {
			return nil, err
		}
		log.Warnf("Retrieved %d with the following error %s", len(result), err)
	}

	return result, nil
}

func ProcessEvents(events []*clogs.FilteredLogEvent) {
	w := csv.NewWriter(os.Stdout)
	var processedEvent []string
	var err error

	w.Write([]string{"Time", "Source", "Caller identity", "Event name", "Parameters", "Error"})
	for _, logEvent := range events {
		if processedEvent, err = ExtractInfoFromEvent(logEvent); err != nil {
			log.Warn(err)
			continue
		}
		w.Write(processedEvent)
	}
	w.Flush()
}

func ExtractInfoFromEvent(logEvent *clogs.FilteredLogEvent) (res []string, err error) {
	var event CloudTrailEvent

	if logEvent.Message == nil {
		return nil, fmt.Errorf("Empty message for event ID: %s", *logEvent.EventId)
	}
	if err = json.Unmarshal([]byte(*logEvent.Message), &event); err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshall log event ID: %s", *logEvent.EventId)
	}

	prettyParams := formatRequestParams(event.RequestParameters)
	identity := formatUsername(event)
	return []string{event.EventTime.Format(time.RFC3339), event.SourceIP, identity, event.EventName, prettyParams, event.ErrorCode}, nil
}

func formatUsername(event CloudTrailEvent) string {
	switch event.UserIdentity.Type {
	case "IAMUser":
		return event.UserIdentity.UserName
	case "Root":
		return "Root"
	case "AWSService":
		return event.UserIdentity.InvokedBy
	case "AssumedRole":
		return userFromARN(event.UserIdentity.Arn)
	default:
		return event.UserIdentity.Principal
	}
}

func userFromARN(input string) string {
	slicedUser := strings.Split(input, "/")
	if len(slicedUser) < 2 {
		return input
	}
	return strings.Join(slicedUser[len(slicedUser)-2:], "/")
}

func formatRequestParams(input interface{}) string {
	prettyParams, _ := json.MarshalIndent(input, "", "")
	out := bytes.ReplaceAll(prettyParams, []byte("\n"), []byte(" "))
	return string(out)
}

type CloudTrailEvent struct {
	EventTime    time.Time `json:"eventTime"`
	EventType    string    `json:"eventType"`
	SourceIP     string    `json:"sourceIPAddress"`
	UserIdentity struct {
		Type        string `json:"type"`
		Principal   string `json:"principalId"`
		UserName    string `json:"userName"`
		Arn         string `json:"arn"`
		AccessKeyID string `json:"accessKeyId"`
		InvokedBy   string `json:"invokedBy"`
	} `json:"userIdentity"`
	EventName         string      `json:"eventName"`
	RequestParameters interface{} `json:"requestParameters"`
	UserAgent         string      `json:"userAgent"`
	EventSource       string      `json:"eventSource"`
	Region            string      `json:"awsRegion"`
	ErrorCode         string      `json:"errorCode"`
}
