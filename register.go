package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Node struct {
	Kind       string   `json:"kind,omitempty"`
	APIVersion string   `json:"apiVersion,omitempty"`
	Metadata   Metadata `json:"metadata,omitempty"`
	Spec       Spec     `json:"spec,omitempty"`
}

type Metadata struct {
	Name string `json:"name,omitempty"`
}

type Spec struct {
	ExternalID string `json:"externalID,omitempty"`
}

type NodeResp struct {
	Reason string `json:"reason,omitempty"`
}

func register(endpoint, addr string) error {
	var n Node
	n.Kind = "Node"
	n.APIVersion = "v1beta3"
	n.Metadata.Name = addr
	n.Spec.ExternalID = addr

	data, err := json.Marshal(&n)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/nodes", endpoint)

	res, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 202 || res.StatusCode == 200 || res.StatusCode == 201 {
		log.Printf("registered machine: %s\n", addr)
		return nil
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 409 {
		return fmt.Errorf("error registering: %s %d %s", addr, res.StatusCode, string(data))
	}

	nr := &NodeResp{}
	if err := json.Unmarshal([]byte(data), &nr); err != nil {
		return err
	}

	if res.StatusCode == 409 && nr.Reason == "AlreadyExists" {
		return nil
	}

	return fmt.Errorf("error registering: %s %s", addr, nr.Reason)
}
