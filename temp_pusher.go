package main

import "encoding/json"

// FeishuData -
type FeishuData struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text"`
}

func pushFeishuMessage(title, text string) {
	d := FeishuData{
		Title: title,
		Text:  text,
	}
	data, err := json.Marshal(&d)
	if err != nil {
		logger.Errorln(err)
		return
	}
	go httpJSON(FeishuCustomBot, map[string]string{"json": string(data)})
}
