package model

type RequestInfo struct {
	Id                   int
	Url                  string
	RequestType          string
	ResponseCode         int64
	ExpectedResponseCode int64
	ResponseTime         int64
	ExpectedResponseTime int64
	ResponseBody         string
	ExpectedResponseBody string
}

type ErrorInfo struct {
	Id           int
	Url          string
	RequestType  string
	ResponseCode int
	ResponseBody string
	Reason       error
	OtherInfo    string
}
