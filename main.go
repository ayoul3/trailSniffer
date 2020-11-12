package main

import (
	"flag"
	"time"

	"github.com/ayoul3/trailSniffer/cloudwatchlogs"
	"github.com/ayoul3/trailSniffer/lib"
	log "github.com/sirupsen/logrus"
)

var logGroup, startDateStr, endDateStr string
var filter *lib.Filter
var startDate, endDate time.Time

func init() {
	filter = &lib.Filter{}
	flag.StringVar(&logGroup, "logGroup", "", "logGroup on Cloudwatch logs")
	flag.StringVar(&filter.AccessKey, "accessKey", "", "Filter on access key")
	flag.StringVar(&filter.UserName, "user", "", `Filter on user name`)
	flag.StringVar(&filter.RoleName, "role", "", `Filter on role name`)
	flag.StringVar(&filter.Account, "account", "", `Filter on the account number`)
	flag.StringVar(&filter.Service, "service", "", `Filter on AWS service name. e.g. ec2, s3, iam, etc.`)
	flag.StringVar(&filter.Event, "event", "", `Filter on event Name: GetObject, etc.`)
	flag.StringVar(&filter.Raw, "filter", "", `Keyword to look for or CloudWatchLog filter, e.g. $.userIdentity.userName="terraform"`)
	flag.StringVar(&startDateStr, "start", "", "Start date to search logs, format YYYY-mm-dd. Defaults to 3 days ago")
	flag.StringVar(&endDateStr, "end", "", "End date to search logs, format YYYY-mm-dd. Defaults to today")
	flag.Parse()
	validateParams()
}

func main() {
	client := cloudwatchlogs.NewClient(cloudwatchlogs.NewAPI())
	res, err := lib.LookupLogs(client, logGroup, filter, startDate, endDate)
	if err != nil {
		log.Fatalf("Got error looking up events: %s", err)
	}
	if len(res) == 0 {
		log.Fatalf("No events matched the filter and date range")
	}
	log.Infof("Got %d events", len(res))
	lib.ProcessEvents(res)
}

func validateParams() {
	var err error
	if logGroup == "" {
		log.Fatalf("Got empty logGroup")
	}
	if startDate, err = time.Parse("2006-01-02", startDateStr); err != nil {
		startDate = time.Now().Add(-3 * 24 * time.Hour)
	}
	if endDate, err = time.Parse("2006-01-02", endDateStr); err != nil {
		endDate = time.Now()
	}
}
