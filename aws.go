package main

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"fmt"
	"encoding/json"
)

type AWSTags struct {
	Tags []AWSTag
}

type AWSTag struct {
	ResourceType string
	ResourceId string
	Value string
	Key string
}

//AWS is a dummy struct
type AWS struct {
	InstanceID string
	Region string
	Tags map[string]string
}

func getInstanceID() string {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return ""

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""

	}
	return string(body)
}

func getRegion() string {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone")
	if err != nil {
		return ""

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""

	}
	az := string(body)
	return az[:len(az) - 1]
}

func getTags(region, instanceID string) map[string]string {
	ret := make(map[string]string)

	args := []string{"ec2", "describe-tags", "--region", region, "--filters", fmt.Sprintf("Name=resource-id,Values=%s", instanceID)}
	cmd := exec.Command("aws",args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error!")
	}

	var awsTags AWSTags
	json.Unmarshal(stdoutStderr, & awsTags)

	for _, tag := range awsTags.Tags {
		ret[tag.Key] = tag.Value
	}

	fmt.Println(ret)

	return ret
}

//Initialize returns an initialized AWS struct
func Initialize() (aws AWS) {
	aws.InstanceID = getInstanceID()
	aws.Region = getRegion()

	aws.Tags = getTags(aws.Region, aws.InstanceID)

	return aws
}
