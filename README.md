# Endpoint Monitor

Monitor any url via this GO application. Get notified through Slack or E-mail when response time is greater than expected, or does not respond as expected.

## Complete Version using InfluxDb

![alt text](https://github.com/patdaman/endpoint-monitor/raw/master/screenshots/graphana.png "Graphana Screenshot")

You can save data to influx db and view response times over a period of time as above using graphana.

[Guide to install influxdb and grafana](https://github.com/patdaman/endpoint-monitor/blob/master/Config.md#database) 

With Endpoint Monitor you can monitor all your REST APIs by adding api details to config file as below.A Notification will be triggered when you api is down or response time is more than expected.

```json
{
	"url":"http://mywebsite.com/v1/data",
	"requestType":"POST",
	"headers":{
		"Authorization":"Bearer ac2168444f4de69c27d6384ea2ccf61a49669be5a2fb037ccc1f",
		"Content-Type":"application/json"
	},
	"formParams":{
		"description":"test",
		"url":"http://google.com"
	},
	"checkEvery":30,
	"responseCode":200,		
	"responseTime":800
},

{
	"url":"http://mywebsite.com/v1/data",
	"requestType":"GET",
	"headers":{
		"Authorization":"Bearer ac2168444f4de69c27d6384ea2ccf61a49669be5a2fb037ccc1f",		
	},
	"urlParams":{
		"name":"endpoint-monitor"
	},
	"checkEvery":300,
	"responseCode":200,		
	"responseTime":800
},

{
	"url":"http://something.com/v1/data",
	"requestType":"DELETE",
	"formParams":{
		"name":"endpoint-monitor"
	},
	"checkEvery":300,
	"responseCode":200,		
	"responseTime":800
}

```
[Guide to write config.json file](https://github.com/patdaman/endpoint-monitor/blob/master/Config.md#writing-a-config-file)

[Sample config.json file](https://github.com/patdaman/endpoint-monitor/blob/master/sample_config.json)

To run the app

```
$ ./endpoint-monitor --config config.json &
```

To run as background process add & at the end

```
$ ./endpoint-monitor --config config.json &	
```
to stop the process 
```
$ jobs
$ kill %jobnumber
```

## Database

Save Requests response time information and error information to your database by adding database details to config file. Currently only Influxdb 0.9.3+ is supported.

## Notifications

Notifications will be triggered when mean response time is below given response time for a request or when an error is occured . Currently the below clients are supported to receive notifications.For more information on setup [click here](https://github.com/patdaman/endpoint-monitor/blob/master/Config.md#notifications)

1. [Slack](https://github.com/patdaman/endpoint-monitor/blob/master/Config.md#slack)
2. [Smtp Email](https://github.com/patdaman/endpoint-monitor/blob/master/Config.md#e-mail)
4. [Http EndPoint](https://github.com/patdaman/endpoint-monitor/blob/master/Config.md#http-endpoint)

Adding support to other clients is simple.[view details](https://github.com/patdaman/endpoint-monitor/blob/master/Config.md#write-your-own-notification-client)

## Contribution

Contributions are welcomed and greatly appreciated. Create an issue if you find bugs.
Send a pull request if you have written a new feature or fixed an issue .Please make sure to write test cases.
