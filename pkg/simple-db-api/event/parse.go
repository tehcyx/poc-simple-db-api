package event

import (
	"github.com/tehcyx/simple-db-api/pkg/store"
)

// GetKymaData transform arbitrary bytes and tries to map this to a Kyma data  struct, to then transform it into StorageData
func GetKymaData(data []byte) (store.StorageData, error) {
	// var newEvent Kyma

	// err := json.Unmarshal(data, &newEvent)
	// if err != nil {
	// 	return store.StorageData{}, fmt.Errorf("Failed to unmarshal json data: %w", err)
	// }

	return store.StorageData{}, nil

}
