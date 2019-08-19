package notify

import (
	"os"
	"testing"

	"github.com/patdaman/endpoint-monitor/model"
)

func TestAddEmptyNotifyObject(t *testing.T) {

	notificationTypes := NotificationTypes{MailNotify{},
		SlackNotify{},
		HttpNotify{},
		OpsgenieNotify{}}

	AddNew(notificationTypes)

	if len(notificationsList) != 0 {
		t.Error("Empty Notification Object should not be added to list")
	}
}

func TestAddValidNotifyObject(t *testing.T) {

	notificationTypes := NotificationTypes{MailNotify{},
		SlackNotify{},
		HttpNotify{"http://statusOk.com", "GET", nil},
		OpsgenieNotify{}}

	AddNew(notificationTypes)

	if len(notificationsList) != 1 {
		t.Error("Failed to Add Notification Object to list")
	}
}

// Send Test model.Notification to all registered clients .To make sure everything is working
func SendTestNotification() {

	println("Sending Test Notifications to the registered clients")

	for _, value := range model.NotificationsList {
		err := value.SendResponseTimeNotification(ResponseTimeNotification{"http://test.com", "GET", 700, 800})

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
