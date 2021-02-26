package main

import (
	"io/ioutil"
	"net/http"
)

//AWS is a dummy struct
type AWS struct {
	InstanceID string
}

//Initialize returns an initialized AWS struct
func Initialize() (aws AWS) {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return aws

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return aws

	}
	aws.InstanceID = string(body)

	return aws
}
