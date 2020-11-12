package lib

import (
	"fmt"
	"strings"
)

type Filter struct {
	Account   string
	AccessKey string
	UserName  string
	RoleName  string
	Service   string
	Event     string
	Raw       string
}

func (f *Filter) Parse() (filter string) {
	filterBuilder := make([]string, 0)
	if f.regularFieldsEmpty() {
		return f.Raw
	}
	if f.Account != "" {
		filterBuilder = append(filterBuilder, fmt.Sprintf(`$.RecipientAccountId = "%s"`, f.Account))
	}
	if f.AccessKey != "" {
		filterBuilder = append(filterBuilder, fmt.Sprintf(`$.userIdentity.accessKeyId = "%s"`, f.AccessKey))
	}
	if f.UserName != "" {
		filterBuilder = append(filterBuilder, fmt.Sprintf(`$.userIdentity.userName = "%s"`, f.UserName))
	}
	if f.RoleName != "" {
		filterBuilder = append(filterBuilder, fmt.Sprintf(`$.userIdentity.arn = "*%s*"`, f.RoleName))
	}
	if f.Service != "" {
		filterBuilder = append(filterBuilder, fmt.Sprintf(`$.eventSource = "%s.amazonaws.com"`, f.Service))
	}
	if f.Event != "" {
		filterBuilder = append(filterBuilder, fmt.Sprintf(`$.eventType = "%s"`, f.Event))
	}
	if f.Raw != "" {
		raw := strings.ReplaceAll(f.Raw, "{", "")
		raw = strings.ReplaceAll(raw, "}", "")
		filterBuilder = append(filterBuilder, raw)
	}

	return fmt.Sprintf("{ %s }", strings.Join(filterBuilder, " && "))

}

func (f *Filter) regularFieldsEmpty() bool {
	return f.Account == "" && f.AccessKey == "" && f.UserName == "" && f.RoleName == "" && f.Service == "" && f.Event == ""
}
