package database

import (
	"errors"
	"testing"
)

func TestInitialize(t *testing.T) {

	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 10, 10)
	if len(responseMean) != len(ids) {
		t.Error("Ids not initialized")
	}

	if MeanResponseCount != 10 {
		t.Error("Mean Response Count Not Set")
	}

	if ErrorCount != 10 {
		t.Error("ErrorCount Not Set")
	}

}

func TestMeanResponseCalculation(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	addResponseTimeToRequest(1, 10)
	mean, _ := getMeanResponseTimeOfUrl(1)

	if mean != 10 {
		t.Error("getMeanResponseTimeOfUrl Failed")
	}

}

func TestAddRequestAndErrorInfo(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	requestInfo := RequestInfo{
		Id:                   1,
		Url:                  "http://test.com",
		RequestType:          "GET",
		ResponseCode:         200,
		ExpectedResponseCode: 200,
		ResponseTime:         10,
		ExpectedResponseTime: 200,
		ResponseBody:         "This is the test body",
		ExpectedResponseBody: "This is the test body",
	}

	errorInfo := ErrorInfo{
		Id:           0,
		Url:          "http://test.com",
		RequestType:  "GET",
		ResponseCode: 0,
		ResponseBody: "this is the body",
		Reason:       errors.New("test error"),
		OtherInfo:    "test other info",
	}

	AddErrorInfo(errorInfo)

	AddRequestInfo(requestInfo)

	mean, err := getMeanResponseTimeOfUrl(1)

	if mean != 10 {

		t.Error("Add Request Info Failed ", mean, err, responseMean[1], MeanResponseCount)
	}

}

func TestClearQueue(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	addResponseTimeToRequest(1, 10)

	clearQueue(1)

	if len(responseMean[1]) != 0 {
		t.Error("ClearQueue Function is not working")
	}
}

func TestAddEmptyDatabase(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	influxDb := InfluxDb{}

	databaseTypes := DatabaseTypes{influxDb}

	AddNew(databaseTypes)

	if len(dbList) != 0 {
		t.Error("Empty Databse should not be added to list")
	}
}

func TestAddValidDatabase(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	influxDb := InfluxDb{"localhost", 8086, "Monitoring", "", ""}

	databaseTypes := DatabaseTypes{influxDb}

	AddNew(databaseTypes)

	if len(dbList) != 1 {
		t.Error("Not able to add databse to list")
	}
}

// Insert test data to database
func addTestErrorAndRequestInfo() {

	println("Adding Test data to your database ....")

	// requestInfo := RequestInfo{0, "http://test.com", "GET", 0, "", 0, 0}
	requestInfo := RequestInfo{
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
	errorInfo := ErrorInfo{
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
