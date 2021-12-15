package alert

import "time"

var Messages []AlertMessage

type AlertMessage struct {

	Status string `json:"status"`

	CommonAnnotations struct {
		Description string `json:"description"`
		Summary string `json:"summary"`
	}`json:"commonAnnotations"`

	Alert []Alerts `json:"alerts"`

	CommonLabels struct{
		Alertname string `json:"alertname"`
		Instance string `json:"instance"`
		Job	string `json:"job"`
		KeyWord string `json:"key_word"`
		Level string `json:"level"`
		Name string `json:"name"`
		Service string `json:"service"`
		Value string `json:"value"`
	}`json:"commonLabels"`
}

type Alerts struct{
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
}