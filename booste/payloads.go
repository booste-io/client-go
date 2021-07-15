package booste

import "encoding/json"

// This package defines the entirety of what data can be accepted on the inbound Post requests

// InStartV1 Defines the acceptable payload into the task/start/v1 endpoint
type inStartV1 struct {
	ID      string        `json:"id" xml:"id" form:"id"`
	Created int64         `json:"created" xml:"created" form:"created"`
	Data    inStartV1Data `json:"data" xml:"data" form:"data"`
}

type inStartV1Data struct {
	APIKey          string      `json:"apiKey" xml:"apiKey" form:"apiKey"`
	ModelKey        string      `json:"modelKey" xml:"modelKey" form:"modelKey"`
	ModelParameters interface{} `json:"modelParameters" xml:"modelParameters" form:"modelParameters"` // Keeps the dynamic substruct as raw bytes
}

// OutStartV1 defines the payload to be returned by the task/start/v1 endpoint
type outStartV1 struct {
	ID         string         `json:"id" xml:"id" form:"id"`
	Message    string         `json:"message" xml:"message" form:"message"`
	Success    bool           `json:"success" xml:"success" form:"success"`
	Created    int64          `json:"created" xml:"created" form:"created"`
	APIVersion string         `json:"apiVersion" xml:"apiVersion" form:"apiVersion"`
	Data       outStartV1Data `json:"data" xml:"data" form:"data"`
}

type outStartV1Data struct {
	TaskID string `json:"taskID" xml:"taskID" form:"taskID"`
}

// InCheckV1 Defines the acceptable payload into the task/check/v1 endpoint
type inCheckV1 struct {
	ID       string        `json:"id" xml:"id" form:"id"`
	Created  int64         `json:"created" xml:"created" form:"created"`
	LongPoll bool          `json:"longPoll" xml:"longPoll" form:"longPoll"`
	Data     inCheckV1Data `json:"data" xml:"data" form:"data"`
}

type inCheckV1Data struct {
	TaskID string `json:"taskID" xml:"taskID" form:"taskID"`
}

// OutCheckV1 defines the payload to be returned by the task/check/v1 endpoint
type outCheckV1 struct {
	ID         string         `json:"id" xml:"id" form:"id"`
	Message    string         `json:"message" xml:"message" form:"message"`
	Success    bool           `json:"success" xml:"success" form:"success"`
	Created    int64          `json:"created" xml:"created" form:"created"`
	APIVersion string         `json:"apiVersion" xml:"apiVersion" form:"apiVersion"`
	Data       outCheckV1Data `json:"data" xml:"data" form:"data"`
}

type outCheckV1Data struct {
	TaskStatus string          `json:"taskStatus" xml:"taskStatus" form:"taskStatus"`
	TaskOut    json.RawMessage `json:"taskOut" xml:"taskOut" form:"taskOut"` // Empty interface allows us to pass dynamic structs
}
