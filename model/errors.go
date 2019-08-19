package model

import (
	"errors"
)

var (
	ErrResposeCode   = errors.New("Response code does not Match expected value")
	ErrResposeBody   = errors.New("Response body does not match expected value")
	ErrTimeout       = errors.New("Request Time out Error")
	ErrCreateRequest = errors.New("Invalid Request Config. Not able to create request")
	ErrDoRequest     = errors.New("Request failed")
)
