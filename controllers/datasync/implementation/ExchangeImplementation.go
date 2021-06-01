package implementation

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"subscription-billing-engine/entities"
	"subscription-billing-engine/entities/exchange"
	"time"
)

type ExchangeProviderDataSync struct {
}

func (request *ExchangeProviderDataSync) IngestByte(byteData []byte) error {
	if err := json.Unmarshal(byteData, &request); err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	return nil
}

func (request *ExchangeProviderDataSync) ProcessDataSync(byteData string) (error, *entities.GenericDataSyncPayload) {
	fmt.Println("request", byteData)
	dataSync := exchange.DataSyncPayload{}
	err := xml.Unmarshal([]byte(byteData), &dataSync)

	fmt.Println("err", err)

	if err != nil {
		return err, nil
	}

	fmt.Println("Exchange Data Sync", dataSync)

	genericResponse := entities.GenericDataSyncPayload{
		ProductId: dataSync.Body.SyncOrderRelation.ProductId,
		SpId:      dataSync.Body.SyncOrderRelation.SpId,
		ServiceId: dataSync.Body.SyncOrderRelation.ServiceID,
		Msisdn:    dataSync.Body.SyncOrderRelation.UserId.ID,
		Meta: dataSync.Body.SyncOrderRelation.ExtensionInfo,
	}

	startTime, err := time.Parse("20060102150405", dataSync.Body.SyncOrderRelation.EffectiveTime)
	if err != nil{
		startTime = time.Now()
	}

	endTime, err := time.Parse("20060102150405", dataSync.Body.SyncOrderRelation.ExpiryTime)
	if err != nil{
		endTime = time.Now()
	}


	genericResponse.StartTime = uint(startTime.Unix())
	genericResponse.EndTime = uint(endTime.Unix())


	if dataSync.Body.SyncOrderRelation.UpdateType == "2"{
		genericResponse.Mode = "unsubscription"
	}else {
		genericResponse.Mode = "subscription"
	}

	for _, v := range dataSync.Body.SyncOrderRelation.ExtensionInfo {
		if v.Key != "transactionID"{
			continue
		}

		genericResponse.Reference = v.Value
	}
	//fmt.Println("Exchange Data Sync", dataSync)

	return nil, &genericResponse
}
