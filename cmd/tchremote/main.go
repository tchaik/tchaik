// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
tchremote is a tool which uses the Tchaik REST API to act as a remote control. See --help for more details.
*/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	HostEnv      = "TCH_ADDR"
	PlayerKeyEnv = "TCH_PLAYER_KEY"
)

var host string
var key string
var keys bool
var action string
var value string

func init() {
	flag.StringVar(&host, "addr", "", fmt.Sprintf("schema://host(:port) address of the REST API (or set %v)", HostEnv))
	flag.StringVar(&key, "key", "", fmt.Sprintf("the key which identifies the player to send actions to (or set %v)", PlayerKeyEnv))
	flag.BoolVar(&keys, "keys", false, "list all the keys on the host")
	flag.StringVar(&action, "action", "", "action to send to the player (requires -key, some require -value)")
	flag.StringVar(&value, "value", "", "value to send to the player")
}

func main() {
	flag.Parse()

	if host == "" {
		host = os.Getenv(HostEnv)
	}
	if host == "" {
		fmt.Printf("must use -addr or set %v\n", HostEnv)
		os.Exit(1)
	}

	if keys {
		list, err := getPlayerKeys()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, k := range list {
			fmt.Println(k)
		}
		return
	}

	if key == "" {
		key = os.Getenv(PlayerKeyEnv)
	}
	if key == "" {
		fmt.Printf("must use -key or set %v\n", PlayerKeyEnv)
		os.Exit(1)
	}

	err := handleAction(action, value)
	if err != nil {
		fmt.Printf("error handling action: %v", err)
		os.Exit(1)
	}
}

func handleAction(action, value string) error {
	if value != "" {
		switch action {
		case "setTime", "setVolume":
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			return sendPlayerAction(f)

		case "setVolumeMute":
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			return sendPlayerAction(b)
		}
	}
	return sendPlayerAction(nil)
}

func getPlayerKeys() ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%v/api/players/", host))
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %v", string(body))
	}

	data := struct {
		Keys []string `json:"keys"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	return data.Keys, nil
}

func sendPlayerAction(value interface{}) error {
	data := struct {
		Action string      `json:"action"`
		Value  interface{} `json:"value,omitempty"`
	}{
		Action: action,
		Value:  value,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling JSON request body: %v", err)
	}
	requestURL := fmt.Sprintf("%v/api/players/%v", host, key)
	req, err := http.NewRequest("PUT", requestURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		return fmt.Errorf("error: %v", strings.TrimSpace(string(body)))
	}
	return nil
}
