package xray

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/noorbala7418/ajor-xui-exporter/internal/model"
	"github.com/sirupsen/logrus"
)

// loginXUI Logins to X-UI panel and returns cookie.
func loginXUI() []*http.Cookie {
	client := &http.Client{}
	loginCred := fmt.Sprintf(`{
		"username" : "%s",
		"password" : "%s",
		"LoginSecret": ""
	}`, os.Getenv("XPANEL_USERNAME"), os.Getenv("XPANEL_PASSWORD"))

	req, err := http.NewRequest("POST", os.Getenv("XPANEL_URL")+"/login", strings.NewReader(loginCred))
	if err != nil {
		logrus.Error("Error create login client. ", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error in login to xui. ", err)
		return nil
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error in login to xui. Error in fetch body.", err)
		return nil
	}
	if resp.StatusCode != 200 {
		logrus.Error("Error in login to XUI. status code is ", resp.StatusCode)
		return nil
	}
	logrus.Debug("Login succeeded. Cookie fetched. Length of cookie is ", len(resp.Cookies()))
	return resp.Cookies()
}

// getInbounds Connects to X-UI panel and takes all inbounds in json mode. It returns list of inbounds.
func getInbounds() (*model.Inbounds, error) {
	loginCookie := loginXUI()[0]
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("XPANEL_URL")+"/xui/inbound/list", nil)
	if err != nil {
		logrus.Error("Error create get data client. ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(loginCookie)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error in send request for get xui inbounds. ", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error in fetch body.", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logrus.Error("Error in get inbounds. status code is ", resp.StatusCode)
		return nil, fmt.Errorf("status code is not 200")
	}

	var inboundList model.Inbounds
	if cleanupInbounds(respBody, &inboundList) != nil {
		logrus.Error("Could not parse inbounds json. ", err)
		return nil, fmt.Errorf("error in parse json %q", err)
	}
	logrus.Debug("Get inbounds success.")
	return &inboundList, nil
}

// cleanupInbounds clears inbounds and merges related informations from inboundSettings with Clients.
func cleanupInbounds(input []byte, inbounds *model.Inbounds) error {
	// Step 1: Get inbounds
	if err := json.Unmarshal(input, &inbounds); err != nil {
		logrus.Error("Could not parse inbounds json. ", err)
		return err
	}

	// Step 2: get client settings and merge them to clinets
	for _, inbound := range inbounds.Inbounds {
		var inboundSettings model.Settings
		if err := json.Unmarshal([]byte(inbound.Settings), &inboundSettings); err != nil {
			logrus.Error("error in unmarshlling settings json", err)
			return err
		}

		for i := 0; i < len(inbound.Clients); i++ {
			for _, clientSetting := range inboundSettings.Clients {
				if inbound.Clients[i].Name == clientSetting.Name {
					inbound.Clients[i].ID = clientSetting.ID
					inbound.Clients[i].AdminEnabled = clientSetting.Enable
					inbound.Clients[i].RemainTraffic = inbound.Clients[i].TotalTraffic - (inbound.Clients[i].DownloadTraffic + inbound.Clients[i].UploadTraffic)
				}
			}
		}
	}

	return nil
}

// func logoutXUI() {}

// GetAllClients returns list of all clients in string slice. You can pass a number to use pagiantion.
func GetAllClients() []model.Client {
	inbounds, _ := getInbounds()

	var users []model.Client

	// Step 1: Get all clients
	for _, inbound := range inbounds.Inbounds {
		for _, client := range inbound.Clients {
			users = append(users, client)
		}
	}

	if len(users) == 0 {
		return nil
	}
	return users
}
