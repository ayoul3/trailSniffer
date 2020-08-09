package lib_test

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	clogs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/ayoul3/trailSniffer/cloudwatchlogs"
	"github.com/ayoul3/trailSniffer/lib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LookupLogs", func() {
	Context("When calling fetchlogs fails", func() {
		It("should return an error", func() {
			client := cloudwatchlogs.NewClient(&cloudwatchlogs.MockAPI{ShouldFail: true})
			_, err := lib.LookupLogs(client, "default", "special-user", time.Now().Add(-3*time.Hour), time.Now())
			Expect(err).To(HaveOccurred())
		})
	})
	Context("When calling fetchlogs returns empty results", func() {
		It("should return empty set", func() {
			client := cloudwatchlogs.NewClient(&cloudwatchlogs.MockAPI{ShouldBeEmpty: true})
			res, err := lib.LookupLogs(client, "default", "special-user", time.Now().Add(-3*time.Hour), time.Now())
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(Equal(0))
		})
	})
	Context("When calling fetchlogs succeeds", func() {
		It("should return 3 events", func() {
			client := cloudwatchlogs.NewClient(&cloudwatchlogs.MockAPI{})
			res, err := lib.LookupLogs(client, "default", "special-user", time.Now().Add(-3*time.Hour), time.Now())
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(Equal(3))
		})
	})
	Context("When calling fetchlogs partially", func() {
		It("should return 3 events", func() {
			client := cloudwatchlogs.NewClient(&cloudwatchlogs.MockAPI{ShouldPartiallyFail: true})
			res, err := lib.LookupLogs(client, "default", "special-user", time.Now().Add(-3*time.Hour), time.Now())
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(Equal(3))
		})
	})
})

var _ = Describe("ExtractInfoFromEvent", func() {
	Context("When the event lacks a message", func() {
		It("should return an error", func() {
			event := &clogs.FilteredLogEvent{EventId: aws.String("123")}
			_, err := lib.ExtractInfoFromEvent(event)
			Expect(err).To(HaveOccurred())
		})
	})
	Context("When the unmarshalling fails", func() {
		It("should return an error", func() {
			message := `"test":"test"`
			event := &clogs.FilteredLogEvent{EventId: aws.String("123"), Message: &message}
			_, err := lib.ExtractInfoFromEvent(event)
			Expect(err).To(HaveOccurred())
		})
	})
	Context("When processing succeeds", func() {
		client := cloudwatchlogs.NewClient(&cloudwatchlogs.MockAPI{})
		events, _ := lib.LookupLogs(client, "default", "special-user", time.Now().Add(-3*time.Hour), time.Now())
		Context("When event type is AWSService", func() {
			It("should succeed", func() {
				processedEvent, err := lib.ExtractInfoFromEvent(events[0])
				Expect(err).ToNot(HaveOccurred())
				Expect(processedEvent[1]).To(Equal("autoscaling.amazonaws.com"))
				Expect(processedEvent[2]).To(Equal("autoscaling.amazonaws.com"))
			})
		})
		Context("When event type is AssumedRole", func() {
			It("should succeed", func() {
				processedEvent, err := lib.ExtractInfoFromEvent(events[1])
				Expect(err).ToNot(HaveOccurred())
				Expect(processedEvent[1]).To(Equal("36.149.154.19"))
				Expect(processedEvent[2]).To(Equal("mytest-role.ec2/i-0faac5b395d439771"))
				Expect(processedEvent[5]).To(Equal("Client.UnauthorizedOperation"))
			})
		})
		Context("When event type is IAMUser", func() {
			It("should succeed", func() {
				processedEvent, err := lib.ExtractInfoFromEvent(events[2])
				Expect(err).ToNot(HaveOccurred())
				Expect(processedEvent[1]).To(Equal("90.154.45.151"))
				Expect(processedEvent[2]).To(Equal("terraform"))
			})
		})
	})
})
