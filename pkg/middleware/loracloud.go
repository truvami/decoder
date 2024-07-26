package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	LoraCloudBaseURL = "https://mgs.loracloud.com/api/v1/device/send"
)

type LoracloudClient struct {
	token string
}

func NewLoraCloudClient(token string) LoracloudClient {
	return LoracloudClient{token}
}

type LoraCloudUplink struct {
	MsgType string `json:"msgtype"`
	Fcnt    int    `json:"fcnt"`
	Port    int    `json:"port"`
	Payload string `json:"payload"`
}

type LoraCloudRequestBody struct {
	DevEui string          `json:"deveui"`
	Uplink LoraCloudUplink `json:"uplink"`
}

func (m LoracloudClient) SolveUplink(uplink LoraCloudRequestBody) error {
	jsonStr, err := json.Marshal(uplink)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", LoraCloudBaseURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", m.token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}
	if err != nil {
		return err
	}

	return nil
}
