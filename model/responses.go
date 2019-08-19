package model

type MessageType string

var (
	// ResponseTimeMessage = "Warning notification From Endpoint Monitor" +
	// 	"\n\nAPI response is returning slower than expected." +
	// 	"\n\nDetails below:" +
	// 	"\nUrl: %v \nRequestType: %v \nCurrent Average Response Time: %v ms\nExpected Response Time: %v ms\n"
	ResponseTimeMessage MessageType = "Warning notification From Endpoint Monitor" +
		"\n\nAPI response is returning slower than expected." +
		"\n\nDetails below:" +
		"\nUrl: %v \nRequestType: %v \nCurrent Average Response Time: %v ms\nExpected Response Time: %v ms\n"
	ResponseCodeMessage MessageType = "ResponseCode"
	ResponseBodyMessage MessageType = "ResponseBody"
	CertMismatchMessage MessageType = "CertMismatch"
	CertExpiringMessage MessageType = "CertExpiring"
	CertExpiredMessage  MessageType = "CertExpired"
	ErrorMessage        MessageType = "Error notification From Endpoint Monitor" +
		"\n\nCannot communicate with " +
		"\n\nPlease find the Details below" +
		"\n\nUrl: %v \nRequestType: %v \nError Message: %v \nOther Info:%v\n"
		// "\n\nUrl: %v \nRequestType: %v \nError Message: %v \nResponse Body: %v\nOther Info:%v\n"+
		// "\n\nThanks", errorNotification.Url, errorNotification.RequestType, errorNotification.Error, errorNotification.ResponseBody, model.ErrorNotification.OtherInfo)
	InformationMessage MessageType = "Information"
)
