package models

import (
	"encoding/json"
	"time"
)

type FunctionCreateResponse struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]interface{} `json:"headers"`
	Body       Function               `json:"body"`
}

type FunctionListResponse struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]interface{} `json:"headers"`
	Body       struct {
		Functions []Function `json:"functions"`
	} `json:"body"`
}

type Function struct {
	CaPort                int32     `json:"caPort"`
	CodeChecksum          string    `json:"codeChecksum"`
	CodeSize              int       `json:"codeSize"`
	Cpu                   float32   `json:"cpu"`
	CreatedTime           time.Time `json:"createdTime"`
	CustomContainerConfig struct {
		AccelerationInfo interface{} `json:"accelerationInfo"`
		AccelerationType string      `json:"accelerationType"`
		Args             string      `json:"args"`
		Command          interface{} `json:"command"`
		Image            string      `json:"image"`
		WebServerMode    bool        `json:"webServerMode"`
	} `json:"customContainerConfig"`
	Description           string            `json:"description"`
	DiskSize              int               `json:"diskSize"`
	EnvironmentVariables  map[string]string `json:"environmentVariables"`
	FunctionId            string            `json:"functionId"`
	FunctionName          string            `json:"functionName"`
	Handler               string            `json:"handler"`
	InitializationTimeout int               `json:"initializationTimeout"`
	InstanceConcurrency   int32             `json:"instanceConcurrency"`
	InstanceType          string            `json:"instanceType"`
	LastModifiedTime      time.Time         `json:"lastModifiedTime"`
	MemorySize            int32             `json:"memorySize"`
	Runtime               string            `json:"runtime"`
	Timeout               int               `json:"timeout"`
}

type TriggerListResponse struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]interface{} `json:"headers"`
	Body       struct {
		Triggers []Trigger `json:"triggers"`
	} `json:"body"`
}

type TriggerCreateResponse struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]interface{} `json:"headers"`
	Body       Trigger                `json:"body"`
}

type Trigger struct {
	TriggerName      string        `json:"triggerName"`
	Description      string        `json:"description"`
	TriggerId        string        `json:"triggerId"`
	TriggerType      string        `json:"triggerType"`
	Qualifier        string        `json:"qualifier"`
	UrlInternet      string        `json:"urlInternet"`
	UrlIntranet      string        `json:"urlIntranet"`
	TriggerConfig    TriggerConfig `json:"triggerConfig"`
	CreatedTime      time.Time     `json:"createdTime"`
	LastModifiedTime time.Time     `json:"lastModifiedTime"`
}

type TriggerConfig struct {
	Methods            []string `json:"methods"`
	AuthType           string   `json:"authType"`
	DisableURLInternet bool     `json:"disableURLInternet"`
}

func (t *TriggerConfig) UnmarshalJSON(b []byte) error {
	// 这玩意是个 string，默认要解析两次
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return err
	}
	t.Methods = make([]string, 0)
	for _, k := range m["methods"].([]interface{}) {
		t.Methods = append(t.Methods, k.(string))
	}
	t.AuthType = m["authType"].(string)
	t.DisableURLInternet = m["disableURLInternet"].(bool)
	return nil
}
