package cloudwatchlogs

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

// MockAPI is an CloudwatchLogs client mock
type MockAPI struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	ShouldFail          bool
	ShouldPartiallyFail bool
	ShouldBeEmpty       bool
}

func (c *MockAPI) FilterLogEventsPages(input *cloudwatchlogs.FilterLogEventsInput, fn func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool) (err error) {
	if c.ShouldBeEmpty {
		return
	}
	if c.ShouldFail {
		return errors.New("failed to fetch events")
	}
	event1 := &cloudwatchlogs.FilteredLogEvent{
		EventId: aws.String("123"),
		Message: aws.String(`{"eventVersion":"1.05","userIdentity":{"type":"AWSService","invokedBy":"autoscaling.amazonaws.com"},"eventTime":"2020-08-08T19:37:33Z","eventSource":"sts.amazonaws.com","eventName":"AssumeRole","awsRegion":"eu-west-1","sourceIPAddress":"autoscaling.amazonaws.com","userAgent":"autoscaling.amazonaws.com","requestParameters":{"roleArn":"arn:aws:iam::466204357409:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling","roleSessionName":"AutoScaling","durationSeconds":1800},"responseElements":{"credentials":{"accessKeyId":"ASIA44ZRK6WSY2EXMAEU","expiration":"Aug 8, 2020 8:07:33 PM","sessionToken":"IQoJb3JpZ2luX2VjEEQaCWV1LXdlc3QtMSJGMEQC"},"assumedRoleUser":{"assumedRoleId":"AROA44ZRK6WS4YIRA6Z5K:AutoScaling","arn":"arn:aws:sts::466204357409:assumed-role/AWSServiceRoleForAutoScaling/AutoScaling"}},"requestID":"9d47ecc9-4759-43a4-b1c0","eventID":"24c9a3cc-72ed-47b8-9dca-01e4ec3566cc","resources":[{"accountId":"466204357409","type":"AWS::IAM::Role","ARN":"arn:aws:iam::466204357409:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling"}],"eventType":"AwsApiCall","recipientAccountId":"466204357409","sharedEventID":"3017954b-eada-403d-a1c1-217029579db2"}`),
	}
	event2 := &cloudwatchlogs.FilteredLogEvent{
		EventId: aws.String("123"),
		Message: aws.String(`{"eventVersion":"1.05","userIdentity":{"type":"AssumedRole","principalId":"AROA44ERK6QSQXKXZDZTQ:i-0faac5b395d439771","arn":"arn:aws:sts::466204357409:assumed-role/mytest-role.ec2/i-0faac5b395d439771","accountId":"466204357409","accessKeyId":"ASIA44ZRK6WSZ6L6LRCD","sessionContext":{"sessionIssuer":{"type":"Role","principalId":"AROA44ERK6QSQXKXZDZTQ","arn":"arn:aws:iam::466204357409:role/mytest-role.ec2","accountId":"466204357409","userName":"mytest-role.ec2"},"webIdFederationData":{},"attributes":{"mfaAuthenticated":"false","creationDate":"2020-08-08T19:29:49Z"}}},"eventTime":"2020-08-08T19:32:39Z","eventSource":"ec2.amazonaws.com","eventName":"DescribeInstances","awsRegion":"eu-west-1","sourceIPAddress":"36.149.154.19","userAgent":"aws-cli/1.18.69 Python/3.6.9 Linux/5.3.0-1032-aws botocore/1.16.19","errorCode":"Client.UnauthorizedOperation","errorMessage":"You are not authorized to perform this operation.","requestParameters":{"instancesSet":{},"filterSet":{}},"responseElements":null,"requestID":"7046a92d-db1e-45e5-8795-7104418b4784","eventID":"bb640ea2-4051-474e-8d49-767fa46bf433","eventType":"AwsApiCall","recipientAccountId":"466204357409"}`),
	}
	event3 := &cloudwatchlogs.FilteredLogEvent{
		EventId: aws.String("123"),
		Message: aws.String(`{"eventVersion":"1.05","userIdentity":{"type":"IAMUser","principalId":"AIDA44ZRK6WS7K6IQQQ3V","arn":"arn:aws:iam::466204357409:user/terraform","accountId":"466204357409","accessKeyId":"AKIA44ZRK6WS4G7MGL6W","userName":"terraform"},"eventTime":"2020-08-08T14:41:36Z","eventSource":"ec2.amazonaws.com","eventName":"DescribeSecurityGroups","awsRegion":"eu-west-1","sourceIPAddress":"90.154.45.151","userAgent":"aws-cli/1.16.301 Python/3.6.8 Linux/4.4.0-18362-Microsoft botocore/1.13.37","requestParameters":{"securityGroupSet":{},"securityGroupIdSet":{"items":[{"groupId":"sg-14d57ebc77858ce2e"}]},"filterSet":{}},"responseElements":null,"requestID":"9b523a82-0c74-46ae-8524-cba381f69636","eventID":"0343474e-6b3a-4570-9265-74acd6c4c83b","eventType":"AwsApiCall","recipientAccountId":"466204357409"}`),
	}

	fn(&cloudwatchlogs.FilterLogEventsOutput{
		Events: []*cloudwatchlogs.FilteredLogEvent{event1, event2, event3},
	}, true)

	if c.ShouldPartiallyFail {
		err = errors.New("failed to fetch some events")
	}

	return
}
