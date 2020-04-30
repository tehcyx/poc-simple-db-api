package event

import "time"

type Kyma struct {
	Date time.Time
	Data string `json:"orderCode"`
}
