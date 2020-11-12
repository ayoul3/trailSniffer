package lib_test

import (
	"github.com/ayoul3/trailSniffer/lib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse", func() {
	Context("When defining multiple filters", func() {
		It("should return a properly formatted filter", func() {

			filter := &lib.Filter{
				UserName:  "user",
				RoleName:  "role",
				Service:   "ec2",
				Account:   "123123123",
				AccessKey: "AKIAAZEAZE",
				Event:     "DescribeInstances",
				Raw:       `{ $.sourceIPAddress = "1.2.3.4" }`,
			}
			filterStr := filter.Parse()
			Expect(filterStr).To(Equal(`{ $.RecipientAccountId = "123123123" && $.userIdentity.accessKeyId = "AKIAAZEAZE" && $.userIdentity.userName = "user" && $.userIdentity.arn = "*role*" && $.eventSource = "ec2.amazonaws.com" && $.eventType = "DescribeInstances" &&  $.sourceIPAddress = "1.2.3.4"  }`))
		})
	})
})
