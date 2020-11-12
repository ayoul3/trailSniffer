package cloudwatchlogs

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

// Client is a CloudwatchLogs custom client
type Client struct {
	api cloudwatchlogsiface.CloudWatchLogsAPI
}

// NewClient instantiates a NewClient - a wrapper over Cloudwatchlogs api
func NewClient(api cloudwatchlogsiface.CloudWatchLogsAPI) *Client {
	return &Client{
		api,
	}
}

// NewAPI returns a new real CloudwatchLogs client
func NewAPI() *cloudwatchlogs.CloudWatchLogs {
	return cloudwatchlogs.New(session.Must(session.NewSession()))
}

func (c *Client) FetchLogs(logGroup, filter string, startDate, endDate time.Time) (output []*cloudwatchlogs.FilteredLogEvent, err error) {
	input := &cloudwatchlogs.FilterLogEventsInput{
		StartTime:     aws.Int64(startDate.Unix() * 1000),
		EndTime:       aws.Int64(endDate.Unix() * 1000),
		Limit:         aws.Int64(10000),
		LogGroupName:  aws.String(logGroup),
		FilterPattern: aws.String(filter),
	}
	err = c.api.FilterLogEventsPages(input,
		func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
			output = append(output, page.Events...)
			return !lastPage
		})
	return output, err
}
