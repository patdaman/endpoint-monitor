package database

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/patdaman/endpoint-monitor/model"
	"github.com/patdaman/endpoint-monitor/notify"
	"github.com/sirupsen/logrus"
)

var (
	MeanResponseCount = 5 // Number of response times to calcuate mean response time
	ErrorCount        = 1 // Threshold to send notification

	dbList       []Database
	responseMean map[int][]int64
	dbMain       Database

	// ErrResposeCode   = errors.New("Response code does not Match expected value")
	// ErrResposeBody   = errors.New("Response body does not match expected value")
	// ErrTimeout       = errors.New("Request Time out Error")
	// ErrCreateRequest = errors.New("Invalid Request Config. Not able to create request")
	// ErrDoRequest     = errors.New("Request failed")

	isLoggingEnabled = false
)

type Database interface {
	Initialize() error
	GetDatabaseName() string
	AddRequestInfo(requestInfo model.RequestInfo) error
	AddErrorInfo(errorInfo model.ErrorInfo) error
}

type DatabaseTypes struct {
	InfluxDb InfluxDb `json:"influxDb"`
}

//Intialize responseMean app and counts
func Initialize(ids map[int]int64, mMeanResponseCount int, mErrorCount int) {

	if mMeanResponseCount != 0 {
		MeanResponseCount = mMeanResponseCount
	}

	if mErrorCount != 0 {
		ErrorCount = mErrorCount
	}
	//TODO: try to make all slices as pointers
	responseMean = make(map[int][]int64)

	for id, _ := range ids {
		queue := make([]int64, 0)
		responseMean[id] = queue
	}

}

// Add database to the database List
func AddNew(databaseTypes DatabaseTypes) {

	v := reflect.ValueOf(databaseTypes)

	for i := 0; i < v.NumField(); i++ {
		dbString := fmt.Sprint(v.Field(i).Interface().(Database))

		// Check whether notify object is empty. If not add to list
		if !isEmptyObject(dbString) {
			dbList = append(dbList, v.Field(i).Interface().(Database))
		}
	}

	if len(dbList) != 0 {
		println("Intializing Database....")
	}

	// Intialize all databases given by user
	for _, value := range dbList {

		initErr := value.Initialize()

		if initErr != nil {
			println("Failed to Intialize Database ")
			os.Exit(3)
		}
	}

	// Set first database as primary
	if len(dbList) != 0 {
		dbMain = dbList[0]
		AddTestErrorAndRequestInfo()
	} else {
		fmt.Println("No Database selected.")
	}
}

// Function called by requests package when request was successfull
// Request data is inserted to all registered databases
func AddRequestInfo(requestInfo model.RequestInfo) {
	logRequestInfo(requestInfo)

	//Insert to all databses
	for _, db := range dbList {
		go db.AddRequestInfo(requestInfo)
	}

	// Response time to queue
	addResponseTimeToRequest(requestInfo.Id, requestInfo.ResponseTime)

	// Calculate current mean response time. If less than expected send notification
	mean, meanErr := getMeanResponseTimeOfUrl(requestInfo.Id)

	if meanErr == nil {
		if mean > requestInfo.ExpectedResponseTime {
			clearQueue(requestInfo.Id)
			// notify.SendResponseTimeNotification(notify.ResponseTimeNotification{
			// 	Url:                  requestInfo.Url,
			// 	RequestType:          requestInfo.RequestType,
			// 	ExpectedResponsetime: requestInfo.ExpectedResponseTime,
			// 	MeanResponseTime:     mean})
			notify.SendNotification(getNotificationObject(requestInfo))
		}
	}
}

// Called by requests package when a request fails
// Error Information inserted to all registered databases
func AddErrorInfo(errorInfo model.ErrorInfo) {
	logErrorInfo(errorInfo)

	// Request failed send notification
	notify.SendErrorNotification(model.ErrorNotification{
		Url:         errorInfo.Url,
		RequestType: errorInfo.RequestType,
		Error:       errorInfo.Reason.Error(),
		OtherInfo:   errorInfo.OtherInfo,
		// ResponseBody: errorInfo.ResponseBody,
	})

	// Add Error information to database
	for _, db := range dbList {
		go db.AddErrorInfo(errorInfo)
	}
}

func addResponseTimeToRequest(id int, responseTime int64) {
	if responseMean != nil {
		queue := responseMean[id]

		if len(queue) == MeanResponseCount {
			queue = queue[1:]
			queue = append(queue, responseTime)
		} else {
			queue = append(queue, responseTime)
		}

		responseMean[id] = queue
	}
}

// Calculate current mean response time for given request id
func getMeanResponseTimeOfUrl(id int) (int64, error) {

	queue := responseMean[id]
	if len(queue) < MeanResponseCount {
		return 0, errors.New("Count has not been reached")
	}

	var sum int64
	for _, val := range queue {
		sum = sum + val
	}

	return sum / int64(MeanResponseCount), nil
}

func clearQueue(id int) {
	responseMean[id] = make([]int64, 0)
}

func isEmptyObject(objectString string) bool {

	objectString = strings.Replace(objectString, "0", "", -1)
	objectString = strings.Replace(objectString, "map", "", -1)
	objectString = strings.Replace(objectString, "[]", "", -1)
	objectString = strings.Replace(objectString, " ", "", -1)

	if len(objectString) > 2 {
		return false
	} else {
		return true
	}
}

func EnableLogging(fileName string) {

	isLoggingEnabled = true

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if len(fileName) == 0 {
		// Output to stderr instead of stdout, could also be a file.
		logrus.SetOutput(os.Stderr)
	} else {
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			println("Invalid File Path given for parameter --log")
			os.Exit(3)
		}

		logrus.SetOutput(f)
	}
}

func logErrorInfo(errorInfo model.ErrorInfo) {

	if isLoggingEnabled {
		logrus.WithFields(logrus.Fields{
			"id":           errorInfo.Id,
			"url":          errorInfo.Url,
			"requestType":  errorInfo.RequestType,
			"responseCode": errorInfo.ResponseCode,
			"reason":       errorInfo.Reason.Error(),
			"otherInfo":    errorInfo.Reason,
			// "responseBody": errorInfo.ResponseBody,
		}).Error("Status Ok Error occurred for url " + errorInfo.Url)
	}

}

func logRequestInfo(requestInfo model.RequestInfo) {

	if isLoggingEnabled {
		logrus.WithFields(logrus.Fields{
			"id":                   requestInfo.Id,
			"url":                  requestInfo.Url,
			"requestType":          requestInfo.RequestType,
			"responseCode":         requestInfo.ResponseCode,
			"responseTime":         requestInfo.ResponseTime,
			"expectedResponseTime": requestInfo.ExpectedResponseTime,
			// "responseBody":         requestInfo.ResponseBody,
		}).Info("")
	}
}

func getNotificationObject(requestInfo model.RequestInfo) model.Notification {
	notificationObject := model.Notification{
		Url:                  requestInfo.Url,
		RequestType:          requestInfo.RequestType,
		ExpectedResponseTime: requestInfo.ExpectedResponseTime,
		ResponseTime:         requestInfo.ResponseTime,
		ExpectedResponseCode: requestInfo.ExpectedResponseCode,
		ResponseCode:         requestInfo.ResponseCode,
	}
	return notificationObject
}

// Insert test data to database
func AddTestErrorAndRequestInfo() {

	println("Adding Test data to your database ....")

	// requestInfo := RequestInfo{0, "http://test.com", "GET", 0, "", 0, 0}
	requestInfo := model.RequestInfo{
		Id:                   0,
		Url:                  "http://test.com",
		RequestType:          "GET",
		ResponseCode:         0,
		ExpectedResponseCode: 200,
		ResponseTime:         0,
		ExpectedResponseTime: 0,
		ResponseBody:         "",
		ExpectedResponseBody: "",
	}

	// errorInfo := ErrorInfo{0, "http://test.com", "GET", 0, "test response", errors.New("test error"), "test other info"}
	errorInfo := model.ErrorInfo{
		Id:           0,
		Url:          "http://test.com",
		RequestType:  "GET",
		ResponseCode: 0,
		Reason:       errors.New("test error"),
		OtherInfo:    "test other info"}

	for _, db := range dbList {
		reqErr := db.AddRequestInfo(requestInfo)
		if reqErr != nil {
			println(db.GetDatabaseName, ": Failed to insert Request Info to database. Please check whether database is installed properly")
		}

		errErr := db.AddErrorInfo(errorInfo)

		if errErr != nil {
			println(db.GetDatabaseName, ": Failed to insert Error Info to database. Please check whether database is installed properly")
		}

	}
}
