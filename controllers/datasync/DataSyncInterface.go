package datasync

import "subscription-billing-engine/entities"

type NetworkDataSyncInterface interface {
	IngestByte([]byte) error
	ProcessDataSync(string) (error, *entities.GenericDataSyncPayload)
}