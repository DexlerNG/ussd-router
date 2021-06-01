package charge

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"subscription-billing-engine/startups/queues"
	"subscription-billing-engine/utils"
)


func QueueChargeHandler(c echo.Context) error {

	//verify if provider is supported
	providersInterface := GetProvider(c.Param("provider"))
	if providersInterface == nil {
		return utils.ValidationResponse(c, "Charge Provider '" + c.Param("provider") + "' is not supported")
	}

	//validate request payload based on provider request data
	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil{
		return utils.ValidationResponse(c, err.Error())
	}
	//convert map string of interface to Class parameter of the interface
	if err = providersInterface.IngestByte(body); err != nil{
		return utils.ValidationResponse(c, err.Error())
	}

	if err = providersInterface.Validate(); err != nil{
		return utils.ValidationResponse(c, err.Error())
	}
	err = queues.RabbitMQPublishToQueue(queues.CHARGE_AUTHORIZATION_PROCESSING_QUEUE, map[string]interface{}{
		"provider": c.Param("provider"),
		"data": providersInterface,
	})

	if  err != nil{
		return utils.ErrorResponse(c, "Unable to Process Payment")
	}

	return utils.AcceptedResponse(c, "Processing...")
}

func AuthorizationCallbackHandler(c echo.Context) error{
	//callbackRequest := map[string]interface{}{}
	//if err := c.Bind(&callbackRequest); err != nil {
	//	return utils.ValidationResponse(c, err.Error())
	//}
	//mv, err := mxj.NewMapXml(string()) // unmarshal

	//verify if provider is supported
	providersInterface := GetProvider(c.Param("provider"))
	if providersInterface == nil {
		return utils.ValidationResponse(c, "Charge Provider '" + c.Param("provider") + "' is not supported")
	}


	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return utils.ValidationResponse(c, err.Error())
	}
	fmt.Println("Body", string(body))
	err = queues.RabbitMQPublishToQueue(queues.CHARGE_PROCESSING_QUEUE, map[string]interface{}{
		"provider": c.Param("provider"),
		"reference": c.Param("reference"),
		"data": string(body),
	})
	if  err != nil{
		return utils.ErrorResponse(c, "Unable to Process Payment")
	}



	//-  Return Specified response
	//-  from the queue
	//-  get the  transactionId
	//-  get data from cache
	//-  initial charge
	//-  get response and process
	return providersInterface.GetAuthorizationResponse(c)
}
//func CreateClientHandler(c echo.Context) error {
//	clientRequest := requests.ClientRequest{}
//	if err := c.Bind(&clientRequest); err != nil {
//		return nil
//	}
//	//update createdBy
//	user := c.Get("user").(models.User)
//	clientRequest.CreatedBy = user.Id
//
//	//validate Request
//	err := clientRequest.ValidateCreateClient()
//	if err != nil {
//		println("Error: ", err)
//		return utils.ValidationResponse(c, err.Error())
//	}
//	client := models.NewClientFromRequest(&clientRequest)
//
//	client.ClientId = models.GenerateClientId(20)
//	//create default user
//	clientRequest.User.ClientId = client.ClientId
//	clientRequest.User.RoleId = "main"
//	//clientRequest.User.Options = map[string]interface{}{
//	//	"verificationType": "email",
//	//	"redirectURL": clientRequest.User.RedirectURL,
//	//}
//	err, createdUser := services.CreateUser(&clientRequest.User)
//	if err != nil {
//		return utils.ErrorResponse(c, err.Error())
//	}
//	client.DefaultUserId = createdUser.Id
//	err, createdClient := client.Save()
//	if err != nil {
//		return utils.ErrorResponse(c, err.Error())
//	}
//
//	clientAndUserMap := map[string]interface{}{
//		"client": createdClient,
//		"user":   createdUser,
//	}
//	go events.ClientCreatedEventHandler(createdClient, clientRequest)
//	//startups.RabbitMQPublishToExchange("receive.exchange", "client.created", createdClient.ClientId, client, "clientId")
//	return utils.CreatedResponse(c, clientAndUserMap)
//}
//
//
//func FetchClientHandler(c echo.Context) error {
//
//	query := map[string]string{}
//	if err := c.Bind(&query); err != nil {
//		return nil
//	}
//	//result := utils.ConvertMapToMongoDBQuery(query)
//	//return c.JSON(http.StatusOK, result)
//	fmt.Println("Query", query)
//	clients := models.FetchPaginatedClients(query)
//	return c.JSON(http.StatusOK, clients)
//}
//
//func FindClientHandler(c echo.Context) error {
//
//	var err error
//	var client interface{}
//	if c.QueryParam("includeSecret") != "" {
//		fmt.Println("Include secret")
//		err, client = models.FindByClientIdWithSecret(c.Param("clientId"))
//	}else{
//		err, client = models.FindByClientId(c.Param("clientId"))
//	}
//
//	if err != nil {
//		fmt.Println("Find Error", err)
//		return utils.CustomErrorResponse(c, "Client Not Found", http.StatusNotFound)
//	}
//
//	//check if there is no query string, return response immediately
//	if len(c.QueryString()) == 0{
//		return utils.SuccessResponse(c, client)
//	}
//
//	fmt.Println("QueryString", c.QueryString())
//
//	clientMap := utils.StructToMap(client)
//	if c.QueryParam("includeUser") != "" {
//		fmt.Println("Include user")
//		err, user := services.FindUser(clientMap["defaultUserId"].(string))
//		fmt.Println("User", err, user)
//		clientMap["user"] = user
//	}
//	return utils.SuccessResponse(c, clientMap)
//}
//
//func FindClientBySubdomainHandler(c echo.Context) error {
//	err, client := models.FindBySubdomain(c.Param("subdomain"))
//	if err != nil {
//		fmt.Println("Find Error", err)
//		return utils.CustomErrorResponse(c, "Client Not Found", http.StatusNotFound)
//	}
//
//	return utils.SuccessResponse(c, map[string]string{
//		"clientId": client.ClientId,
//		"logo": client.Logo,
//		"subdomain": client.SubDomain,
//		"name": client.Name,
//	})
//}
//
//func CountClientHandler(c echo.Context) error {
//	//query := c.QueryParams().Get("group")
//	//queryParams := c.QueryParams()
//	//fmt.Println("queryParams",queryParams)
//	query := map[string]interface{}{}
//	//if err := c.Bind(&query); err != nil {
//	//	return nil
//	//}
//	//
//	//if query["dateFrom"] != nil {
//	//	fmt.Println("query[dateFrom]", query["dateFrom"])
//	//	//timeConversion, timeError := time.Parse("2018-10-02", query["dateFrom"])
//	//	//fmt.Println("Time Errpr", timeError)
//	//	query["createdAt"] = bson.M{"$lte": time.Now()}
//	//}
//	fmt.Println("Query", query)
//	clientCount := models.CountAllClient(&query)
//	return utils.SuccessResponse(c, clientCount)
//}
//func UpdateClientHandler() {
//
//}
//
//func UpdateAccess(c echo.Context) error {
//	//query := c.QueryParams().Get("group")
//
//	accessType := []string{"granted","revoked"}
//	requestBody := map[string]interface{}{}
//	if err := c.Bind(&requestBody); err != nil {
//		return nil
//	}
//
//	requestBody["access"] = strings.ToLower(requestBody["access"].(string))
//	if !utils.MapContainsString(accessType, requestBody["access"].(string)){
//		return utils.ValidationResponse(c,  "Invalid Access Type")
//	}
//	err, message := models.UpdateOneClient(c.Param("clientId"), requestBody)
//	if err != nil{
//		return utils.ErrorResponse(c, err.Error())
//	}
//	return utils.UpdatedResponse(c, message)
//}
////
////func Delete(c echo.Context) error {
////	//query := c.QueryParams().Get("group")
////	println("ScheduleId", c.Param("scheduleId"))
////	err, id := models.DeleteOneSchedule(c.Param("scheduleId"))
////	if err != nil{
////		return utils.ErrorResponse(c, err.Error())
////	}
////	cron.Reload()
////	return utils.CreatedResponse(c, id)
////
