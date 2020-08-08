package cloudwatchlogs

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/prometheus/common/log"
)

// Client is a EC2 custom client
type Client struct {
	api cloudwatchlogsiface.CloudWatchLogsAPI
}

// NewClient returns a new Client from an CloudWatch client
func NewClient(api cloudwatchlogsiface.CloudWatchLogsAPI) *Client {
	return &Client{
		api,
	}
}

// NewAPI returns a new concrete CloudwatchLogs client
func NewAPI() *cloudwatchlogs.CloudWatchLogs {
	return cloudwatchlogs.New(session.Must(session.NewSession()))
}

func (c *Client) FetchLogs(logGroup, filter string, startDate, endDate time.Time) (output []*cloudwatchlogs.FilteredLogEvent, err error) {

	preparedFilter := prepareFilter(filter)
	input := &cloudwatchlogs.FilterLogEventsInput{
		StartTime:     aws.Int64(startDate.Unix() * 1000),
		EndTime:       aws.Int64(endDate.Unix() * 1000),
		Limit:         aws.Int64(10000),
		LogGroupName:  aws.String(logGroup),
		FilterPattern: aws.String(preparedFilter),
	}
	log.Infof("Full filter %s", preparedFilter)
	err = c.api.FilterLogEventsPages(input,
		func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
			log.Infof("Got %d events", len(page.Events))
			output = append(output, page.Events...)
			return !lastPage
		})
	return output, err
}

func prepareFilter(input string) string {
	if strings.HasPrefix(input, "$") {
		return fmt.Sprintf("{ %s }", input)
	}
	return input
}