package notify

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/patdaman/endpoint-monitor/src/model"
)

type SamanageAPIHeader string

func (h SamanageAPIHeader) String() string {
	return string(h)
}

const (
	SamanageUrl              SamanageAPIHeader = "https://api.samanage.com/"
	SamanageTokenKey         SamanageAPIHeader = "X-Samanage-Authorization"
	SamanageAcceptAPIKey     SamanageAPIHeader = "Accept"
	SamanageAcceptAPIValue   SamanageAPIHeader = "application/vnd.samanage.v2.1+json"
	SamanageContentTypeKey   SamanageAPIHeader = "Content-Type"
	SamanageContentTypeValue SamanageAPIHeader = "application/json"
)

type SamanageNotify struct {
	RequestType        string            `json:"requestType"`
	Headers            map[string]string `json:"headers"`
	SamanageTokenValue SamanageAPIHeader `json:"token"`
}

func (samanageNotify SamanageNotify) GetClientName() string {
	return "Notify Samanage"
}

func (samanageNotify SamanageNotify) Initialize() error {
	return nil
}

func (samanageNotify SamanageNotify) SendResponseTimeNotification(responseTimeNotification model.ResponseTimeNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := MessageParam{getMessageFromResponseTimeNotification(responseTimeNotification)}

	if samanageNotify.Headers[ContentType] == JsonContentType {

		jsonBody, jsonErr := GetJsonParamsBody(msgParam)
		if jsonErr != nil {
			return jsonErr
		}
		request, reqErr = http.NewRequest(samanageNotify.RequestType,
			SamanageUrl.String(),
			jsonBody)
	} else {
		urlParams := GetUrlValues(msgParam)
		request, reqErr = http.NewRequest(samanageNotify.RequestType,
			SamanageUrl.String(),
			bytes.NewBufferString(urlParams.Encode()))

		request.Header.Add(ContentType, FormContentType)
		request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
	}

	if reqErr != nil {
		return reqErr
	}

	AddHeaders(request, samanageNotify.Headers)

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http response Status code expected: %v Got : %v ", http.StatusOK, getResponse.StatusCode))
	}

	return nil

}

func (samanageNotify SamanageNotify) SendErrorNotification(errorNotification model.ErrorNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := MessageParam{getMessageFromErrorNotification(errorNotification)}

	if samanageNotify.Headers[ContentType] == JsonContentType {

		jsonBody, jsonErr := GetJsonParamsBody(msgParam)
		if jsonErr != nil {
			return jsonErr
		}
		request, reqErr = http.NewRequest(samanageNotify.RequestType,
			SamanageUrl.String(),
			jsonBody)

	} else if samanageNotify.Headers[ContentType] == FormContentType {
		urlParams := GetUrlValues(msgParam)
		request, reqErr = http.NewRequest(samanageNotify.RequestType,
			SamanageUrl.String(),
			bytes.NewBufferString(urlParams.Encode()))
		request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
	} else {
		urlParams := GetUrlValues(msgParam)
		request, reqErr = http.NewRequest(samanageNotify.RequestType,
			SamanageUrl.String(),
			bytes.NewBufferString(urlParams.Encode()))

		request.Header.Add(ContentType, FormContentType)
		request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
	}

	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}

	AddHeaders(request, samanageNotify.Headers)

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		fmt.Println(respErr, samanageNotify)
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http response Status code expected: %v Got : %v ", http.StatusOK, getResponse.StatusCode))
	}

	return nil
}

// func AddHeaders(req *http.Request, headers map[string]string) {
// 	for key, value := range headers {
// 		req.Header.Add(key, value)
// 	}
// }

// func GetUrlValues(msgParam MessageParam) url.Values {
// 	urlParams := url.Values{}
// 	urlParams.Set("message", msgParam.Message)
// 	return urlParams
// }

// func GetJsonParamsBody(msgParam MessageParam) (io.Reader, error) {

// 	data, jsonErr := json.Marshal(msgParam)

// 	if jsonErr != nil {

// 		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

// 		return nil, jsonErr
// 	}

// 	return bytes.NewBuffer(data), nil
// }

// func getStringFromResponseBody(body io.ReadCloser) string {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(body)
// 	return buf.String()
// }
