package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/patdaman/endpoint-monitor/dbclient"
	cb "github.com/patdaman/endpoint-monitor/requests"
	"github.com/patdaman/endpoint-monitor/tracing"
	"github.com/spf13/viper"
)

var appName = "accountservice"

func init() {
	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configServerURL := flag.String("configServerUrl", "http://configserver:8888", "Address to config server")
	configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")

	flag.Parse()

	viper.Set("profile", *profile)
	viper.Set("configServerURL", *configServerURL)
	viper.Set("configBranch", *configBranch)
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Infof("Starting %v\n", appName)

	config.LoadConfigurationFromBranch(
		viper.GetString("configServerURL"),
		appName,
		viper.GetString("profile"),
		viper.GetString("configBranch"))

	initializeBoltClient()
	initializeMessaging()
	initializeTracing()
	cb.ConfigureHystrix([]string{"imageservice", "quotes-service"}, service.MessagingClient)

	handleSigterm(func() {
		cb.Deregister(service.MessagingClient)
		service.MessagingClient.Close()
	})
	service.StartWebServer(viper.GetString("server_port"))
}
func initializeTracing() {
	tracing.InitTracing(viper.GetString("zipkin_server_url"), appName)
}

func initializeMessaging() {
	if !viper.IsSet("amqp_server_url") {
		panic("No 'amqp_server_url' set in configuration, cannot start")
	}

	service.MessagingClient = &messaging.AmqpClient{}
	service.MessagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))
	service.MessagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
