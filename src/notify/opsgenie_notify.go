package notify

import (
	"errors"
	"fmt"

	"github.com/opsgenie/opsgenie-go-sdk/alertsv2"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"
	"github.com/patdaman/endpoint-monitor/src/model"
)

type OpsgenieNotify struct {
	APIKey string `json:"apikey"`
	Alias  string `json:"alias"`
}

func (opsgenieNotify OpsgenieNotify) GetClientName() string {
	return "OpsGenie Alert"
}

func (opsgenieNotify OpsgenieNotify) Initialize() error {
	return nil
}

func (opsgenieNotify OpsgenieNotify) SendNotification(notification model.Notification) error {
	cli := new(ogcli.OpsGenieClient)
	cli.SetAPIKey(opsgenieNotify.APIKey)

	alertCli, _ := cli.AlertV2()

	request := alertsv2.CreateAlertRequest{
		Message:     notification.Message,
		Alias:       opsgenieNotify.Alias,
		Description: notification.Description,
		Tags:        notification.Tags,
		Details:     notification.Details,
		Entity:      notification.Environment,
		Priority:    Priority(notification.Priority),
		Note:        notification.Note,
	}

	response, err := alertCli.Create(request)
	if err != nil {
		// return err.Error()
		return err
	}
	if len(response.RequestID) == 0 {
		return errors.New(fmt.Sprintf("OpsGenie did not return a new alert ID."))
	}
	return nil
}

func Priority(priority string) alertsv2.Priority {
	switch priority {
	case "Critical", "P1":
		return alertsv2.P1
	case "High", "P2":
		return alertsv2.P2
	case "Low", "P4":
		return alertsv2.P4
	case "Information", "P5":
		return alertsv2.P5
	}
	return alertsv2.P3
}
