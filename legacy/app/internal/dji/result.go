package dji

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	key Key
	sequenceNumber uint32
	value interface{}
	errorCode int64
	errorDescription string
}

func (r Result) ErrorDescription() string {
	return r.errorDescription
}

func (r Result) ErrorCode() int64 {
	return r.errorCode
}

func (r Result) Value() interface{} {
	return r.value
}

func (r Result) SequenceNumber() uint32 {
	return r.sequenceNumber
}

func (r Result) Key() Key {
	return r.key
}

type jsonResult struct {
	Error int64
	Key int
	Tag uint32
	Value map[string]interface{}
}

func NewResultFromJSON(jsonData []byte) *Result {
	if len(jsonData) == 0 {
		return &Result{
			errorCode: -1,
			errorDescription: "empty or nil json data",
		}
	}

	var result jsonResult
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return &Result{
			errorCode: -1,
			errorDescription: fmt.Sprintf("invalid json data: %s", string(jsonData)),
		}
	}

	var value interface{}
	errorDescription := ""
	errorCode := result.Error
	if errorCode != 0 {
		errorDescription = "result error"
	} else {
		ok := false
		value, ok = result.Value["value"]
		if !ok {
			return &Result{
				errorCode: -1,
				errorDescription: "no result value found: %s",
			}
		}
	}

	return &Result{
		keyByValue(result.Key),
		result.Tag,
		value,
		errorCode,
		errorDescription,
	}
}
