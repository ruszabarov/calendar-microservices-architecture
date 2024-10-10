package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func GetMeetingsByIds(ids []string) []Meeting {
	idsParam := strings.Join(ids, ",")

	url := fmt.Sprintf("http://krakend:8080/meetings?ids=%s", idsParam)

	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	var data []Meeting
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Decode Error:", err)
		return nil
	}

	fmt.Printf("Received Data: %+v\n", data)

	return data
}
