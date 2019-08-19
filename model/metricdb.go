package model

type MetricDb struct {
	databaseType struct {
		host         string
		port         int64
		databaseName string
		username     string
		password     string
	}
}
