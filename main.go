package main

import (
	"flag"
	"time"

	"github.com/ayoul3/trailSniffer/cloudwatchlogs"
	"github.com/ayoul3/trailSniffer/lib"
	log "github.com/sirupsen/logrus"
)

var logGroup, filter, startDateStr, endDateStr string
var startDate, endDate time.Time

func init() {
	flag.StringVar(&logGroup, "logGroup", "", "logGroup on Cloudwatch logs")
	flag.StringVar(&filter, "filter", "", `Keyword to look for or Cloudwatchlog filter, e.g. $.userIdentity.userName="terraform"`)
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
