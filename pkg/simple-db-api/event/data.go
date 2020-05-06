package event

import "time"

type Kyma struct {
	Date        time.Time
	BaseSiteUID string `json:"baseSiteUid"`
	OrderCode   string `json:"orderCode"`
}
