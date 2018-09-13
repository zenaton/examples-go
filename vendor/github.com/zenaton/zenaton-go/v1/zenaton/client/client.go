package client

import (
	"fmt"
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
	ZENATON_API_URL     = "https://zenaton.com/api/v1"
	ZENATON_WORKER_URL  = "http://localhost"
	DEFAULT_WORKER_PORT = 4001
	WORKER_API_VERSION  = "v_newton"

	MAX_ID_SIZE = 256

	APP_ENV   = "app_env"
	APP_ID    = "app_id"
	API_TOKEN = "api_token"

	ATTR_ID        = "custom_id"
	ATTR_NAME      = "name"
	ATTR_CANONICAL = "canonical_name"
	ATTR_DATA      = "data"
	ATTR_PROG      = "programming_language"
	ATTR_MODE      = "mode"

	PROG = "Go"

	EVENT_INPUT = "event_input"
	EVENT_NAME  = "event_name"

	WORKFLOW_KILL  = "kill"
	WORKFLOW_PAUSE = "pause"
	WORKFLOW_RUN   = "run"
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
	os.Setenv("ZENATON_LIBRARY_PATH", zenatonDirectory)

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
	workerURL := os.Getenv("ZENATON_WORKER_URL")
	if workerURL == "" {
		workerURL = ZENATON_WORKER_URL
	}

	workerPort := os.Getenv("ZENATON_WORKER_PORT")
	if workerPort == "" {
		workerPort = strconv.Itoa(DEFAULT_WORKER_PORT)
	}

	url := workerURL + ":" + workerPort + "/api/" + WORKER_API_VERSION +
		"/" + resources + "?"

	return c.addAppEnv(url, params)
}

func (c *Client) getWebsiteURL(resources, params string) string {
	apiURL := ZENATON_API_URL
	if os.Getenv("ZENATON_API_URL") != "" {
		apiURL = os.Getenv("ZENATON_API_URL")
	}
	var url = apiURL + "/" + resources + "?" + API_TOKEN + "=" + apiToken + "&"
	return c.addAppEnv(url, params)
}

//todo: figure out how to handle errors
func (c *Client) StartWorkflow(flowName, flowCanonical, customID string, data interface{}) error {

	if len(customID) >= MAX_ID_SIZE {
		//todo: handle this error better
		fmt.Println(`Provided id must not exceed ` + strconv.Itoa(MAX_ID_SIZE) + ` bytes`)
	}

	body := make(map[string]interface{})
	body[ATTR_PROG] = PROG
	body[ATTR_CANONICAL] = flowCanonical
	if flowCanonical == "" {
		body[ATTR_CANONICAL] = nil
	}
	body[ATTR_NAME] = flowName

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

	body[ATTR_DATA] = encodedData
	body[ATTR_ID] = customID

	fmt.Println("c.getInstanceWorkerUrl(), body", c.getInstanceWorkerUrl(""), body)
	resp, err := service.Post(c.getInstanceWorkerUrl(""), body)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			panic("connection refused: try starting zenaton with 'zenaton start'")
		}
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	//todo: panic? or return error?
	if err != nil {
		panic(err)
	}
	//todo: this is ugly
	if strings.Contains(string(respBody), `Your worker does not listen to app`) {
		panic(string(respBody))
	}
	//todo: fix this
	return nil
}

//todo: figure out how to handle errors
func (c *Client) StartWorkflow2(flowName, flowCanonical, customID string, data interface{}) error {

	if len(customID) >= MAX_ID_SIZE {
		//todo: handle this error better
		fmt.Println(`Provided id must not exceed ` + strconv.Itoa(MAX_ID_SIZE) + ` bytes`)
	}

	body := make(map[string]interface{})
	body[ATTR_PROG] = PROG
	body[ATTR_CANONICAL] = flowCanonical
	body[ATTR_NAME] = flowName

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

	body[ATTR_DATA] = encodedData
	body[ATTR_ID] = customID

	resp, err := service.Post(c.getInstanceWorkerUrl(""), body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	//todo: panic? or return error?
	if err != nil {
		panic(err)
	}
	//todo: this is ugly
	if strings.HasPrefix(string(respBody), `{"error":`) {
		panic(string(respBody))
	}
	//todo: fix this
	return nil
}

func (c *Client) KillWorkflow(workflowName, customId string) error {
	err := c.updateInstance(workflowName, customId, WORKFLOW_KILL)
	if err != nil {
		return errors.New(fmt.Sprint("unable to kill workflow: ", workflowName, " error: ", err.Error()))
	}
	return nil
}

func (c *Client) PauseWorkflow(workflowName, customId string) error {
	err := c.updateInstance(workflowName, customId, WORKFLOW_PAUSE)
	if err != nil {
		return errors.New(fmt.Sprint("unable to pause workflow: ", workflowName, " error: ", err.Error()))
	}
	return nil
}

func (c *Client) ResumeWorkflow(workflowName, customId string) error {
	err := c.updateInstance(workflowName, customId, WORKFLOW_RUN)
	if err != nil {
		return errors.New(fmt.Sprint("unable to resume workflow: ", workflowName, " error: ", err.Error()))
	}
	return nil
}

func (c *Client) FindWorkflowInstance(workflowName, customId string) (map[string]map[string]string, error) {
	params := ATTR_ID + "=" + customId + "&" + ATTR_NAME + "=" + workflowName + "&" + ATTR_PROG + "=" + PROG

	resp, err := service.Get(c.getInstanceWebsiteURL(params))
	if err != nil {
		fmt.Println("1")
		return nil, errors.New("unable to find workflow with id: " + customId + " error: " + err.Error())
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if strings.HasPrefix(string(respBody), `{"error":`) {
		return nil, errors.New(string(respBody))
	}

	if err != nil {
		fmt.Println("2")
		return nil, errors.New("unable to find workflow with id: " + customId + " error: " + err.Error())
	}

	var respMap map[string]map[string]string
	err = json.Unmarshal(respBody, &respMap)
	if err != nil {
		//fmt.Println("3", string(respBody))
		return nil, errors.New("unable to find workflow with id: " + customId + " error: " + err.Error())
	}

	return respMap, nil
}

// todo: should this return something?
func (c *Client) SendEvent(workflowName, customID, eventName string, eventData interface{}) {
	var url = c.getSendEventURL()
	body := make(map[string]interface{})
	body[ATTR_PROG] = PROG
	body[ATTR_NAME] = workflowName
	body[ATTR_ID] = customID
	body[EVENT_NAME] = eventName
	encodedData, err := serializer.Encode(eventData)
	if err != nil {
		panic(err)
	}

	if encodedData == "null" {
		encodedData = "{}"
	}
	body[EVENT_INPUT] = encodedData

	fmt.Println("sendEvent body: ", body)
	service.Post(url, body)
}

func (c *Client) updateInstance(workflowName, customId, mode string) error {
	var params = ATTR_ID + "=" + customId
	var body = make(map[string]interface{})
	body[ATTR_PROG] = PROG
	body[ATTR_NAME] = workflowName
	body[ATTR_MODE] = mode
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
