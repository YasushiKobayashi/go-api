package handler

import (
	"app/config"
	"bytes"
	"encoding/json"
	"net/http"

	"golang.org/x/exp/utf8string"
)

type (
	Slack struct {
		Channel  string `json:"channel"`
		Username string `json:"username"`
		IconUrl  string `json:"icon_url"`
		Text     string `json:"text"`
	}
)

func SendSlack(param string) error {
	slack := Slack{}
	slack = Slack{
		Channel:  config.SLACK_CNANNEL,
		Username: config.SLACK_USERNAME,
		IconUrl:  config.SLACK_FACEICON,
		Text:     param,
	}

	json, err := json.Marshal(slack)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		config.SLACK_WEBHOOKURL,
		bytes.NewBuffer([]byte(json)),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

func TrimStr(str string, count int) string {
	cont := utf8string.NewString(str)
	var contCount int = cont.RuneCount()
	if count > contCount {
		count = contCount
	}
	contStr := ">>>" + cont.Slice(0, count)
	return contStr
}
