package services

import (
	"net/http"
)

var req = &http.Request{
	Header: map[string][]string{
		"Content-Type": {"application/json; charset=UTF-8"},
		"source":       {"internal"},
	},
}

//func CreateUser(user *requests.UserRequest) (error, *models.User) {
//	var userBaseURL = os.Getenv("USER_SERVICE_URL")
//	jsonValue, _ := json.Marshal(user)
//	reqURL, _ := url.Parse(userBaseURL + "/v1/users")
//	req.URL = reqURL
//	req.Method = "POST"
//	req.Header.Set("client-id", user.ClientId)
//	req.Body = ioutil.NopCloser(strings.NewReader(string(jsonValue)))
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		// handle error
//		println("User Error: ", err.Error())
//		return err, nil
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		//call slack and drop Message
//		println("Body Reading Slack Error: ", err.Error())
//		return err, nil
//	}
//	println("Response Body: ", string(body))
//
//	userResponse := requests.ServicesStandardResponse{}
//	err = json.Unmarshal(body, &userResponse)
//	if err != nil {
//		fmt.Println("Error While decoding client: ", err)
//		return err, nil
//	}
//	if !utils.IsStringEmpty(userResponse.Error) {
//		return errors.New(userResponse.Error), nil
//	}
//	jsonValue, _ = json.Marshal(userResponse.Data["user"])
//	createdUser := models.User{}
//	err = json.Unmarshal(jsonValue, &createdUser)
//	if err != nil {
//		fmt.Println("Error While decoding user: ", err)
//		return err, nil
//	}
//	return nil, &createdUser
//}
//
//func FindUser(userId string) (error, *models.User) {
//	var userBaseURL = os.Getenv("USER_SERVICE_URL")
//
//	reqURL, _ := url.Parse(userBaseURL + "/v1/users/" + userId)
//	req.URL = reqURL
//	req.Method = "GET"
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		// handle error
//		println("User Error: ", err.Error())
//		return err, nil
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		//call slack and drop Message
//		println("Body Reading Slack Error: ", err.Error())
//		return err, nil
//	}
//	println("Response Body: ", string(body))
//
//	userResponse := requests.ServicesStandardResponse{}
//	err = json.Unmarshal(body, &userResponse)
//	if err != nil {
//		fmt.Println("Error While decoding client: ", err)
//		return err, nil
//	}
//	if !utils.IsStringEmpty(userResponse.Error) {
//		return errors.New(userResponse.Error), nil
//	}
//	jsonValue, _ := json.Marshal(userResponse.Data)
//	createdUser := models.User{}
//	err = json.Unmarshal(jsonValue, &createdUser)
//	if err != nil {
//		fmt.Println("Error While decoding user: ", err)
//		return err, nil
//	}
//	return nil, &createdUser
//}
