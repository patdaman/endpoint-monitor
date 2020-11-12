package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	// "github.com/gorilla/mux"
	"github.com/patdaman/endpoint-monitor/src/database"
	"github.com/patdaman/endpoint-monitor/src/notify"
	"github.com/patdaman/endpoint-monitor/src/requests"

	// "github.com/urfave/cli"
	// "database"
	// "notify"
	// "requests"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/urfave/cli"
)

type configParser struct {
	NotifyWhen    NotifyWhen               `json:"notifyWhen"`
	Requests      []requests.RequestConfig `json:"requests"`
	Notifications notify.NotificationTypes `json:"notifications"`
	Database      database.DatabaseTypes   `json:"database"`
	Concurrency   int                      `json:"concurrency"`
	Port          int                      `json:"port"`
}

type NotifyWhen struct {
	MeanResponseCount int `json:"meanResponseCount"`
	ErrorCount        int `json:"errorCount"`
}

func main() {

	// Cli tool setup to get config file path from parameters
	app := cli.NewApp()
	app.Name = "Endpoint-Monitor"
	app.Usage = "Monitor http(s) endpoints \r\nSave telemetry data to Influx DB \r\nSend notifications in multiple formats"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Value: "config.json",
			Usage: "location of config file",
		},
		&cli.StringFlag{
			Name:  "log",
			Value: "",
			Usage: "file to save logs",
		},
		&cli.BoolFlag{
			Name:  "notify",
			Usage: "send notifications on alerts",
		},
		&cli.BoolFlag{
			Name:  "test",
			Usage: "send test notifications on start",
		},
	}

	app.Action = func(c *cli.Context) error {

		if fileExists(c.String("config")) {

			if len(c.String("log")) != 0 {
				// Log parameter given. Confirm file can be created at given path

				if !logFilePathValid(c.String("log")) {
					println("Invalid File Path given for parameter --log")
					os.Exit(3)
				}
			}

			println("Reading File :", c.String("config"))

			// Start monitoring when a valid file path is given
			startMonitoring(c.String("config"), c.String("log"), c.Bool("notifyActive"), c.Bool("test"))
		} else {
			println("Config file not present at the given location: ", c.String("config"), "\nPlease give correct file location using --config parameter")
		}
		return nil
	}

	// Run as cli app
	app.Run(os.Args)
}

func startMonitoring(configFileName string, logFileName string, notifyActive bool, test bool) {

	configFile, err := os.Open(configFileName)

	if err != nil {
		fmt.Println("Error opening config file:\n", err.Error())
	}

	// Parse the config file data to configParser struct
	jsonParser := json.NewDecoder(configFile)
	var config configParser
	if err = jsonParser.Decode(&config); err != nil {
		fmt.Println("Error parsing config file .Please check format of the file \nParse Error:", err.Error())
		os.Exit(3)
	}

	if notifyActive {
		// Setup notification clients
		notify.AddNew(config.Notifications)

		if test == true {
			// Send test notifications to all notification clients
			notify.SendTestNotification()
		}
	}
	// Create unique ids for each request date given in config file
	reqs, ids := validateAndCreateIdsForRequests(config.Requests)

	// Set up and initialize databases
	database.AddNew(config.Database)
	database.Initialize(ids, config.NotifyWhen.MeanResponseCount, config.NotifyWhen.ErrorCount)

	// Initialize and start monitoring all the apis
	requests.RequestsInit(reqs, config.Concurrency)
	requests.StartMonitoring()

	database.EnableLogging(logFileName)

	// Just to check Endpoint Monitor is running or not
	http.HandleFunc("/", statusHandler)

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	if config.Port == 0 {
		// Default port
		http.ListenAndServe(":7321", r)
	} else {
		// Port is mentioned in config file
		http.ListenAndServe(":"+strconv.Itoa(config.Port), r)
	}
}

// Prints program running
func statusHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Endpoint Monitor is running \n Planning to display useful information in further releases")
}

// Returns whether a file exits or not
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func logFilePathValid(name string) bool {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		return false
	}

	return true
}

// Checks whether each request in config file has valid data
// Creates unique ids for each request using math/rand
func validateAndCreateIdsForRequests(reqs []requests.RequestConfig) ([]requests.RequestConfig, map[int]int64) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Array of ids used by database pacakge to calculate mean response time and send notifications
	ids := make(map[int]int64, 0)

	// Array of new requests data after updating the ids
	newreqs := make([]requests.RequestConfig, 0)

	for i, requestConfig := range reqs {
		validateErr := requestConfig.Validate()
		if validateErr != nil {
			println("\nInvalid Request data in config file for Request #", i, " ", requestConfig.Url)
			println("Error:", validateErr.Error())
			os.Exit(3)
		}

		// Set a random value as id
		randInt := random.Intn(1000000)
		ids[randInt] = requestConfig.ResponseTime
		requestConfig.SetId(randInt)
		newreqs = append(newreqs, requestConfig)
	}

	return newreqs, ids
}
