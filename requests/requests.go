package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/patdaman/endpoint-monitor/database"
)

var (
	RequestsList   []RequestConfig
	requestChannel chan RequestConfig
	throttle       chan int
)

const (
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
	FormContentType = "application/x-www-form-urlencoded"
	JsonContentType = "application/json"

	DefaultTime         = "300s"
	DefaultResponseCode = http.StatusOK
	DefaultConcurrency  = 1

	AuthenticationRequired = http.StatusNetworkAuthenticationRequired
	Unauthorized           = http.StatusUnauthorized
)

type RequestConfig struct {
	Id           int
	Url          string            `json:"url"`
	Environment  string            `json:"environment"`
	Priority     string            `json:"priority"`
	RequestType  string            `json:"requestType"`
	Headers      map[string]string `json:"headers"`
	FormParams   map[string]string `json:"formParams"`
	UrlParams    map[string]string `json:"urlParams"`
	Tags         []string          `json:"tags"`
	Details      map[string]string `json:"details"`
	ResponseCode int               `json:"responseCode"`
	ResponseTime int64             `json:"responseTime"`
	ResponseBody string            `json:"responseBody"`
	CheckEvery   time.Duration     `json:"checkEvery"`
}

// Set Id for request
func (requestConfig *RequestConfig) SetId(id int) {
	requestConfig.Id = id
}

// Check all request Config fields are valid
func (requestConfig *RequestConfig) Validate() error {

	if len(requestConfig.Url) == 0 {
		return errors.New("Invalid Url")
	}

	if _, err := url.Parse(requestConfig.Url); err != nil {
		return errors.New("Invalid Url")
	}

	if len(requestConfig.RequestType) == 0 {
		return errors.New("RequestType cannot be empty")
	}

	if requestConfig.ResponseTime == 0 {
		return errors.New("ResponseTime cannot be empty")
	}

	if requestConfig.ResponseCode == 0 {
		requestConfig.ResponseCode = DefaultResponseCode
	}

	if requestConfig.CheckEvery == 0 {
		defTime, _ := time.ParseDuration(DefaultTime)
		requestConfig.CheckEvery = defTime
	}

	return nil
}

// Initialize data from config file and check all requests
func RequestsInit(data []RequestConfig, concurrency int) {
	RequestsList = data

	// Throttle channel is used to limit number of concurrent requests
	if concurrency == 0 {
		throttle = make(chan int, DefaultConcurrency)
	} else {
		throttle = make(chan int, concurrency)
	}

	requestChannel = make(chan RequestConfig, len(data))

	if len(data) == 0 {
		println("\nNo requests to monitor. Please add requests to your config file")
		os.Exit(3)
	}
	// Send requests to confirm all requests are valid
	println("\nSending requests to apis.....making sure endpoints are reachable and configuration is correct before monitoring")
	println("Api Count: ", len(data))

	for i, requestConfig := range data {
		println("Request #", i, " : ", requestConfig.RequestType, " ", requestConfig.Url)

		// Perform request
		reqErr := PerformRequest(requestConfig, nil)

		if reqErr != nil {
			// Request Failed
			println("\nFailed !!!! Not able to perfome below request")
			println("\n----Request Deatails---")
			println("Url :", requestConfig.Url)
			println("Type :", requestConfig.RequestType)
			println("Error Reason :", reqErr.Error())
			println("\nPlease check the config file and try again")

			os.Exit(3)
		}
	}

	println("All requests Successfull")
}

// Start monitoring by calling createTicker method for each request
func StartMonitoring() {
	fmt.Println("\nStarted Monitoring all ", len(RequestsList), " apis .....")

	go listenToRequestChannel()

	for _, requestConfig := range RequestsList {
		go createTicker(requestConfig)
	}
}

// Time ticker writes data to request channel for every request
func createTicker(requestConfig RequestConfig) {

	var ticker *time.Ticker = time.NewTicker(requestConfig.CheckEvery * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			requestChannel <- requestConfig
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// All tickers write to request channel
func listenToRequestChannel() {

	// Throttle limits number of concurrent requests
	for {
		select {
		case requect := <-requestChannel:
			throttle <- 1
			go PerformRequest(requect, throttle)
		}
	}

}

// Uses date from request Config and creates and executes http request
func PerformRequest(requestConfig RequestConfig, throttle chan int) error {
	// Remove value from throttel channel when request is completed
	defer func() {
		if throttle != nil {
			<-throttle
		}
	}()

	var request *http.Request
	var reqErr error

	if len(requestConfig.FormParams) == 0 {
		// FormParams create a request
		request, reqErr = http.NewRequest(requestConfig.RequestType,
			requestConfig.Url,
			nil)

	} else {
		if requestConfig.Headers[ContentType] == JsonContentType {
			jsonBody, jsonErr := GetJsonParamsBody(requestConfig.FormParams)
			if jsonErr != nil {
				go database.AddErrorInfo(database.ErrorInfo{
					Id:           requestConfig.Id,
					Url:          requestConfig.Url,
					RequestType:  requestConfig.RequestType,
					ResponseCode: 0,
					Reason:       database.ErrCreateRequest,
					OtherInfo:    jsonErr.Error(),
					// ResponseBody: "Unable to create Request object",
				})

				return jsonErr
			}
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				jsonBody)

		} else {
			formParams := GetUrlValues(requestConfig.FormParams)

			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				bytes.NewBufferString(formParams.Encode()))

			request.Header.Add(ContentLength, strconv.Itoa(len(formParams.Encode())))

			if requestConfig.Headers[ContentType] != "" {
				// Add content type to header if user doesnt mention it config file
				// Default content type application/x-www-form-urlencoded
				request.Header.Add(ContentType, FormContentType)
			}
		}
	}

	if reqErr != nil {
		go database.AddErrorInfo(database.ErrorInfo{
			Id:           requestConfig.Id,
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: 0,
			Reason:       database.ErrCreateRequest,
			OtherInfo:    reqErr.Error(),
			// ResponseBody: "Unable to create Request object",
		})

		return reqErr
	}

	// Add url parameters to query if present
	if len(requestConfig.UrlParams) != 0 {
		urlParams := GetUrlValues(requestConfig.UrlParams)
		request.URL.RawQuery = urlParams.Encode()
	}

	// Add headers to the request
	AddHeaders(request, requestConfig.Headers)

	//TODO: put timeout ?
	/*
		timeout := 10 * requestConfig.ResponseTime

		client := &http.Client{
			Timeout: timeout,
		}
	*/

	client := &http.Client{}
	start := time.Now()

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		var statusCode int
		if getResponse == nil {
			statusCode = 0
		} else {
			statusCode = getResponse.StatusCode
		}
		go database.AddErrorInfo(database.ErrorInfo{
			Id:           requestConfig.Id,
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: statusCode,
			Reason:       database.ErrDoRequest,
			OtherInfo:    respErr.Error(),
			// ResponseBody: convertResponseToString(getResponse),
		})
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != requestConfig.ResponseCode {
		go database.AddErrorInfo(database.ErrorInfo{
			Id:           requestConfig.Id,
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: getResponse.StatusCode,
			Reason:       errResposeCode(getResponse.StatusCode, requestConfig.ResponseCode),
			OtherInfo:    "",
			// ResponseBody: convertResponseToString(getResponse),
		})
		return errResposeCode(getResponse.StatusCode, requestConfig.ResponseCode)
	}

	elapsed := time.Since(start)

	// Request succesfull. Add infomartion to Database
	go database.AddRequestInfo(database.RequestInfo{
		Id:                   requestConfig.Id,
		Url:                  requestConfig.Url,
		RequestType:          requestConfig.RequestType,
		ResponseCode:         getResponse.StatusCode,
		ResponseTime:         elapsed.Nanoseconds() / 1000000,
		ExpectedResponseTime: requestConfig.ResponseTime,
		// ResponseBody:         convertResponseToString(getResponse),
	})

	return nil
}

//convert response body to string
func convertResponseToString(resp *http.Response) string {
	if resp == nil {
		return " "
	}
	buf := new(bytes.Buffer)
	_, bufErr := buf.ReadFrom(resp.Body)

	if bufErr != nil {
		return " "
	}

	return buf.String()
}

//Add header values from map to request
func AddHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

//convert params in map to url.Values
func GetUrlValues(params map[string]string) url.Values {
	urlParams := url.Values{}
	i := 0
	for key, value := range params {
		if i == 0 {
			urlParams.Set(key, value)
		} else {
			urlParams.Add(key, value)
		}
	}

	return urlParams
}

//Creates body for request of type application/json from map
func GetJsonParamsBody(params map[string]string) (io.Reader, error) {
	data, jsonErr := json.Marshal(params)

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}

func GetCertificateInfo(w http.ResponseWriter, r *http.Request) {

	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		cn := strings.ToLower(r.TLS.PeerCertificates[0].Subject.CommonName)
		fmt.Printf("CN: %s", cn)
	}

}

//creates an error when response code from server is not equal to response code mentioned in config file
func errResposeCode(status int, expectedStatus int) error {
	return errors.New(fmt.Sprintf("Got Response code %v. Expected Response Code %v ", status, expectedStatus))
}
