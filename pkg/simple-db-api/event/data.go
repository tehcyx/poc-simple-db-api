package event

import "time"

type Kyma struct {
	Date       time.Time
	BaseSiteID string `json:"baseSiteId"`
	OrderCode  string `json:"orderCode"`
}
