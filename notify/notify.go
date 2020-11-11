package notify

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/patdaman/endpoint-monitor/model"
)

// Diffrent types of clients to deliver model.Notifications
type NotificationTypes struct {
	MailNotify MailNotify     `json:"mail"`
	Slack      SlackNotify    `json:"slack"`
	Http       HttpNotify     `json:"httpEndPoint"`
	Opsgenie   OpsgenieNotify `json:"opsgenie"`
	Samanage   SamanageNotify `json:"samanage"`
}

var (
	errorCount        = 0
	NotificationsList []Notify
)

type Notify interface {
	GetClientName() string
	Initialize() error
	SendNotification(notification model.Notification) error

	SendResponseTimeNotification(notification model.ResponseTimeNotification) error
	SendResponseCodeNotification(notification model.ResponseCodeNotification) error
	SendResponseBodyNotification(notification model.ResponseBodyNotification) error
	SendErrorNotification(notification model.ErrorNotification) error
}

// Add model.Notification clients given by user in config file to model.NotificationsList
func AddNew(notificationTypes NotificationTypes) {

	v := reflect.ValueOf(notificationTypes)

	for i := 0; i < v.NumField(); i++ {
		notifyString := fmt.Sprint(v.Field(i).Interface().(Notify))
		// Check whether notify object is empty . if its not empty add to the list
		if !isEmptyObject(notifyString) {
			NotificationsList = append(NotificationsList, v.Field(i).Interface().(Notify))
		}
	}

	if len(NotificationsList) == 0 {
		println("No clients Registered for Notifications")
	} else {
		println("Initializing Notification Clients....")
	}

	for _, value := range NotificationsList {
		initErr := value.Initialize()

		if initErr != nil {
			println("Notifications : Failed to Initialize ", value.GetClientName(), ". Please check the details in config file ")
			println("Error Details : ", initErr.Error())
		} else {
			println("Notifications : ", value.GetClientName(), " Intialized")
		}

	}
}

// Format model.Notification type and send
func SendNotification(notification model.Notification) {
	for _, value := range NotificationsList {
		err := value.SendNotification(notification)
		//TODO: Combine the other two functions into this method
		// 	+ add additional notification types
		//TODO: retry if fails ? what to do when error occurs ?
		if err != nil {

		}
	}
}

// Send response time model.Notification to all clients registered
func SendResponseTimeNotification(responseTimeNotification model.ResponseTimeNotification) {

	for _, value := range NotificationsList {
		err := value.SendResponseTimeNotification(responseTimeNotification)
		//TODO: retry if fails ? what to do when error occurs ?
		if err != nil {

		}
	}
}

// Send Error model.Notification to all clients registered
func SendErrorNotification(errorNotification model.ErrorNotification) {

	for _, value := range NotificationsList {
		err := value.SendErrorNotification(errorNotification)
		//TODO: retry if fails ? what to do when error occurs ?
		if err != nil {

		}
	}
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
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

// A readable message string from ResponseTimeNotification
func getMessageFromResponseTimeNotification(responseTimeNotification model.ResponseTimeNotification) string {
	message := fmt.Sprintf(string(model.ResponseTimeMessage), responseTimeNotification.Url, responseTimeNotification.RequestType, responseTimeNotification.MeanResponseTime, responseTimeNotification.ExpectedResponsetime)
	return message
}

// A readable message string from model.ErrorNotification
func getMessageFromErrorNotification(errorNotification model.ErrorNotification) string {
	message := fmt.Sprintf(string(model.ErrorMessage), errorNotification.Url, errorNotification.RequestType, errorNotification.Error, errorNotification.OtherInfo)
	// message := fmt.Sprintf(model.ErrorMessage, errorNotification.Url, errorNotification.RequestType, errorNotification.Error, errorNotification.ResponseBody, model.ErrorNotification.OtherInfo)
	return message
}

// A readable subject string from model.Notification
func getSubjectFromErrorNotification(notification model.Notification) string {
	subject := fmt.Sprintf("Error returned from %v", notification.Url)
	return subject
}

// A readable subject string from model.Notification
func getSubjectFromNotification(notification model.Notification) string {
	subject := fmt.Sprintf("Error returned from %v", notification.Url)
	return subject
}

// A readable message string from model.ErrorNotification
func getMessageFromNotification(notification model.Notification) string {
	subject := fmt.Sprintf("Error returned from %v", notification.Url)
	return subject
}

// Send Test model.Notification to all registered clients .To make sure everything is working
func SendTestNotification() {

	println("Sending Test Notifications to the registered clients")

	for _, value := range NotificationsList {
		err := value.SendResponseTimeNotification(model.ResponseTimeNotification{"http://test.com", "GET", 700, 800})

		if err != nil {
			println("Failed to Send Response Time Notification to ", value.GetClientName(), " Please check the details entered in the config file")
			println("Error Details :", err.Error())
			os.Exit(3)
		} else {
			println("Sent Test Response Time Notification to ", value.GetClientName(), ". Make sure you received it")
		}

		// err1 := value.SendErrorNotification(model.ErrorNotification{"http://test.com", "GET", "This is test Notification", "Test Notification", "test"})
		err1 := value.SendErrorNotification(model.ErrorNotification{"http://test.com", "GET", "This is test Notification", "Test Notification"})

		if err1 != nil {
			println("Failed to Send Error Notification to ", value.GetClientName(), " Please check the details entered in the config file")
			println("Error Details :", err1.Error())
			os.Exit(3)
		} else {
			println("Sent Test Error Notification to ", value.GetClientName(), ". Make sure you received it")
		}
	}
}
