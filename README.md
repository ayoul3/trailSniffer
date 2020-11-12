# trailSniffer

![Build](https://github.com/ayoul3/trailSniffer/workflows/Go/badge.svg)

A tiny wrapper around the AWS Cloudwatch Logs API that will look up events and format them in CSV.
Very useful when searching for all AWS API calls made by a user, role, IP and so on.
Quickly gather the least number of privileges to assign to a user or role.

## Usage
1. [Prepare a Go environment](https://golang.org/dl/)
2. Make sure you have your AWS credentials set up (don't forget to specify the AWS_REGION as well)
3. Launch the tool. Assuming your Cloudwatch Log group is CloudTrail/DefaultLogGroup
```sh
go run main.go -logGroup CloudTrail/DefaultLogGroup

Time, Source, Caller identity, Event name, Parameters,Error
2020-08-08T19:32:40Z, 34.46.94.19, mytest-role.ec2/i-0faac5b745ea39647, DescribeInstances, "{ ""filterSet"": {}, ""instancesSet"": {} }",
2020-08-08T14:50:57Z, ec2.amazonaws.com, sts.amazonaws.com, AssumeRole, "{ ""roleArn"": ""arn:aws:iam::466204357409:role/mytest-role.ec2"", ""roleSessionName"": ""i-0faac5b895db39671"" }",
2020-08-08T14:41:36Z, 212.105.54.11, special-user, DescribeSecurityGroups,"{ ""filterSet"": {}, ""securityGroupIdSet"": { ""items"": [ { ""groupId"": ""sg-0015278c79864125"" } ] }, ""securityGroupSet"": {} }",
```

TrailSniffer will limit the search to the last three days. You can expand that search with the flags `start` and `end`.

Finally, there are multiple flags to help quickly locate the info you're looking for:

* `go run main.go -username "terraform" -logGroup CloudTrail/DefaultLogGroup`
* `go run main.go -accesskey 'AKIATESTKEYAWS' -logGroup CloudTrail/DefaultLogGroup`
* `go run main.go -event DescribeInstances -logGroup CloudTrail/DefaultLogGroup`

`-raw` flag supports [Cloudwatch filters](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html)

* `go run main.go -raw '($.eventName = "AssumeRole" && $.userAgent = "go*") || $.userIdentity.accountId = "123456789"' -logGroup CloudTrail/DefaultLogGroup`