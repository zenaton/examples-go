package client

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"io/ioutil"

	"errors"
	"path"

	"strings"

	"encoding/json"

	"github.com/zenaton/zenaton-go/v1/zenaton/service"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

const (
	zenatonAPIurl     = "https://zenaton.com/api/v1"
	zenatonWorkerURL  = "http://localhost"
	defaultWorkerPort = 4001
	workerAPIversion  = "v_newton"

	maxIDsize = 256

	APP_ENV   = "app_env"
	APP_ID    = "app_id"
	API_TOKEN = "api_token"

	attrID        = "custom_id"
	attrName      = "name"
	attrCanonical = "canonical_name"
	attrData      = "data"
	attrProg      = "programming_language"
	attrMode      = "mode"

	prog = "Go"

	eventInput = "event_input"
	eventName  = "event_name"

	workflowKill  = "kill"
	workflowPause = "pause"
	workflowRun   = "run"
)

var (
	clientInstance *Client
	appID          string
	apiToken       string
	appEnv         string
)

type Client struct{}

func InitClient(appIDx, apiTokenx, appEnvx string) {
	appID = appIDx
	apiToken = apiTokenx
	appEnv = appEnvx
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	directory := path.Dir(filename)
	zenatonDirectory := directory[:len(directory)-len("/client")]
	err := os.Setenv("ZENATON_LIBRARY_PATH", zenatonDirectory)
	if err != nil {
		panic(err)
	}

}

func NewClient(worker bool) *Client {
	if clientInstance != nil {
		if !worker && (appID == "" || apiToken == "" || appEnv == "") {
			//todo: produce error?
			panic("Please initialize your Zenaton instance with your credentials")
			// throw new ExternalZenatonError('Please initialize your Zenaton instance with your credentials')
		}
		return clientInstance
	}
	return &Client{}
}

func (c *Client) GetWorkerUrl(resources string, params string) string {
	workerURL := os.Getenv("zenatonWorkerURL")
	if workerURL == "" {
		workerURL = zenatonWorkerURL
	}

	workerPort := os.Getenv("ZENATON_WORKER_PORT")
	if workerPort == "" {
		workerPort = strconv.Itoa(defaultWorkerPort)
	}

	url := workerURL + ":" + workerPort + "/api/" + workerAPIversion +
		"/" + resources + "?"

	return c.addAppEnv(url, params)
}

func (c *Client) getWebsiteURL(resources, params string) string {
	apiURL := zenatonAPIurl
	if os.Getenv("ZENATON_API_URL") != "" {
		apiURL = os.Getenv("ZENATON_API_URL")
	}
	var url = apiURL + "/" + resources + "?" + API_TOKEN + "=" + apiToken + "&"
	return c.addAppEnv(url, params)
}

func (c *Client) StartWorkflow(flowName, flowCanonical, customID string, data interface{}) {

	if len(customID) >= maxIDsize {
		panic(`Provided id must not exceed ` + strconv.Itoa(maxIDsize) + ` bytes`)
	}

	body := make(map[string]interface{})
	body[attrProg] = prog
	body[attrCanonical] = flowCanonical
	if flowCanonical == "" {
		body[attrCanonical] = nil
	}
	body[attrName] = flowName

	var encodedData string
	var err error

	if data == nil {
		encodedData = "{}"
	} else {
		encodedData, err = serializer.Encode(data)
		if err != nil {
			panic(err)
		}
	}

	body[attrData] = encodedData
	body[attrID] = customID

	resp, err := service.Post(c.getInstanceWorkerUrl(""), body)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			panic("connection refused: try starting zenaton with 'zenaton start'")
		}
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	if strings.Contains(string(respBody), `Your worker does not listen to app`) {
		panic(string(respBody))
	}
}

func (c *Client) KillWorkflow(workflowName, customId string) error {
	err := c.updateInstance(workflowName, customId, workflowKill)
	if err != nil {
		return errors.New(fmt.Sprint("unable to kill workflow: ", workflowName, " error: ", err.Error()))
	}
	return nil
}

func (c *Client) PauseWorkflow(workflowName, customId string) error {
	err := c.updateInstance(workflowName, customId, workflowPause)
	if err != nil {
		return errors.New(fmt.Sprint("unable to pause workflow: ", workflowName, " error: ", err.Error()))
	}
	return nil
}

func (c *Client) ResumeWorkflow(workflowName, customId string) error {
	err := c.updateInstance(workflowName, customId, workflowRun)
	if err != nil {
		return errors.New(fmt.Sprint("unable to resume workflow: ", workflowName, " error: ", err.Error()))
	}
	return nil
}

func (c *Client) FindWorkflowInstance(workflowName, customId string) (map[string]map[string]string, bool, error) {
	params := attrID + "=" + customId + "&" + attrName + "=" + workflowName + "&" + attrProg + "=" + prog

	resp, err := service.Get(c.getInstanceWebsiteURL(params))
	if err != nil {
		return nil, false, errors.New("1unable to find workflow with id: " + customId + " error: " + err.Error())
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, errors.New("2unable to find workflow with id: " + customId + " error: " + err.Error())
	}

	var respMap map[string]map[string]string

	err = json.Unmarshal(respBody, &respMap)
	if err != nil {
		fmt.Println("respBody: ", resp.StatusCode)
		return nil, false, errors.New("3unable to find workflow with id: " + customId + " error: " + err.Error())
	}

	return respMap, true, nil
}

// todo: should this return something?
func (c *Client) SendEvent(workflowName, customID, name string, eventData interface{}) {
	var url = c.getSendEventURL()
	body := make(map[string]interface{})
	body[attrProg] = prog
	body[attrName] = workflowName
	body[attrID] = customID
	body[eventName] = name
	encodedData, err := serializer.Encode(eventData)
	if err != nil {
		panic(err)
	}

	if encodedData == "null" {
		encodedData = "{}"
	}
	body[eventInput] = encodedData

	service.Post(url, body)
}

func (c *Client) updateInstance(workflowName, customId, mode string) error {
	var params = attrID + "=" + customId
	var body = make(map[string]interface{})
	body[attrProg] = prog
	body[attrName] = workflowName
	body[attrMode] = mode
	_, err := service.Put(c.getInstanceWorkerUrl(params), body)
	return err
}

func (c *Client) getSendEventURL() string {
	return c.GetWorkerUrl("events", "")
}

func (c *Client) getInstanceWebsiteURL(params string) string {
	return c.getWebsiteURL("instances", params)
}

func (c *Client) getInstanceWorkerUrl(params string) string {
	return c.GetWorkerUrl("instances", params)
}

func (c *Client) addAppEnv(url, params string) string {

	var appEnvx string
	if appEnv != "" {
		appEnvx = APP_ENV + "=" + appEnv + "&"
	}

	var appIDx string
	if appID != "" {
		appIDx = APP_ID + "=" + appID + "&"
	}

	if params != "" {
		params = params + "&"
	}

	return url + appEnvx + appIDx + params
}
