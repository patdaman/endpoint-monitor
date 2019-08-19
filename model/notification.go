package model

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/buildkite/interpolate"
)

type Notification struct {
	isActive             bool
	Url                  string
	MessageType          MessageType
	Endpoints            []string
	RequestType          string
	Environment          string
	Priority             string
	ExpectedResponseCode int64
	ResponseCode         int64
	ExpectedResponseTime int64
	ExpectedResponseBody string
	ResponseTime         int64
	ResponseBody         string
	Subject              string
	Description          string
	Message              string
	Error                string
	OtherInfo            string
	Tags                 []string
	Details              map[string]string
	Note                 string
}

func CreateNotification(n Notification) Notification {

	// var notificationVars []string
	// n1 := reflect.ValueOf(n)
	// elems := reflect.ValueOf(&n)Elem()
	// elemNames := elems.Type()
	// notify := make([]interface{}, n1.NumField())

	// for i := 0; i < n1.NumField(); i++ {
	// 	if n1.Field(i).CanInterface() {
	// 		notificationVars = append(notificationVars, elemNames.Field(i).Name+"="+n1.Field(i).Interface())
	// 	}
	// }
	// env := interpolate.NewSliceEnv(notificationVars)
	env := interpolate.NewSliceEnv([]string{
		"URL=" + n.Url,
		"RequestType=" + n.RequestType,
		"Environment=" + n.Environment,
		"Priority=" + n.Priority,
		"ExpectedResponseCode=" + strconv.Itoa(int(n.ExpectedResponseCode)),
		"ResponseCode=" + strconv.Itoa(int(n.ResponseCode)),
		"ExpectedResponseTime=" + strconv.Itoa(int(n.ExpectedResponseTime)),
		"ExpectedResponseBody=" + n.ExpectedResponseBody,
		"ResponseTime=" + strconv.Itoa(int(n.ResponseTime)),
		"ResponseBody=" + n.ResponseBody,
		"Description=" + n.Description,
		"Error=" + n.Error,
		"OtherInfo=" + n.OtherInfo,
		"Tags=" + strings.Join(n.Tags, ", "),
		"Details=" + CreateKeyValuePairs(n.Details),
		"Note=" + n.Note,
	})
	messageTemplate := ""
	subjectTemplate := ""
	n.Message, _ = interpolate.Interpolate(env, messageTemplate)
	n.Subject, _ = interpolate.Interpolate(env, subjectTemplate)
	return n
}

func CreateKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

type ResponseTimeNotification struct {
	Url                  string
	RequestType          string
	ExpectedResponsetime int64
	MeanResponseTime     int64
}

type ErrorNotification struct {
	Url         string
	RequestType string
	Error       string
	OtherInfo   string
	// ResponseBody string
}

type ResponseCodeNotification struct {
	Url                  string
	RequestType          string
	ExpectedResponseCode int64
	ResponseCode         int64
}

type ResponseBodyNotification struct {
	Url                  string
	RequestType          string
	Environment          string
	Priority             string
	ExpectedResponseCode int64
	ResponseCode         int64
	ExpectedResponseBody string
	ResponseBody         string
}
