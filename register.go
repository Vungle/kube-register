package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
  "io/ioutil"
	"log"
	"net/http"
)

type Node struct {
	Kind   string `json:"kind,omitempty"`
	UID     string `json:"uid,omitempty"`
  APIVersion string `json:"apiVersion,omitempty"`
}

type NodeStatus struct {
  HostIP string `json:"hostIP,omitempty"`
}

type NodeResp struct {
  Reason string `json:"reason,omitempty"`
}

func register(endpoint, addr string) error {
	node := &Node{
		Kind: "Node",
    APIVersion: "v1beta3",
		UID: addr,
		&NodeStatus{
			HostIP: addr,
		},
	}
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/v1beta3/nodes", endpoint)
	res, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 202 || res.StatusCode == 201 || res.StatusCode == 200 {
		log.Printf("registered machine: %s\n", addr)
		return nil
	}
	nodeResp := &NodeResp{}
	data, err = ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(data), &nodeResp)
	if res.StatusCode == 409 && nodeResp.Reason == "AlreadyExists" {
		log.Printf("Already registered machine: %s\n", addr)
		return nil
	}
	log.Printf("Response: %#v", res)
	log.Printf("Response Body:\n%s", string(data))
	body, err := ioutil.ReadAll(res.Body)
	reason := ""
	if err == nil {
		reason = ": " + string(body)
	}
	return errors.New("error registering: " + addr + reason)
}
