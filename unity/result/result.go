package result

import (
	"encoding/json"
	"fmt"

	"github.com/brunoga/unitybridge/unity/key"
)

// Reesult represents a result from an operation on a key. The zero value
// is not valid.
type Result struct {
	key       *key.Key
	tag       uint64
	errorCode int32
	errorDesc string
	value     any
}

type jsonResultValue struct {
	Value any `json:"value"`
}

type jsonResult struct {
	Key   uint32
	Tag   uint64
	Error int32
	Value jsonResultValue
}

// NewFromJSON creates a new Result from the given JSON data. Any errors are
// reported in the Result itself and should be handled by anyone that cares
// about it.
func NewFromJSON(jsonData []byte) *Result {
	r := &Result{}

	if len(jsonData) == 0 {
		r.errorCode = -1
		r.errorDesc = "empty or nil json data"
		return r
	}

	jr := jsonResult{}

	err := json.Unmarshal(jsonData, &jr)
	if err != nil {
		r.errorCode = -1
		r.errorDesc = fmt.Sprintf("error unmarshalling json data: %s",
			err.Error())
		return r
	}

	key, err := key.FromSubType(jr.Key)
	if err != nil {
		r.errorCode = -1
		r.errorDesc = fmt.Sprintf("error creating key from sub type %d: %s",
			jr.Key, err.Error())
		return r
	}

	errorDesc := "no error"
	if jr.Error != 0 {
		errorDesc = fmt.Sprintf("error %d", jr.Error)
	}

	// TODO(bga): Make sure all values actually have a "value" field.
	value := jr.Value.Value

	r.key = key
	r.tag = jr.Tag
	r.errorCode = jr.Error
	r.errorDesc = errorDesc
	r.value = value

	return r
}

// Key returns the key associated with this result.
func (r *Result) Key() *key.Key {
	return r.key
}

// Tag returns the tag associated with this result.
func (r *Result) Tag() uint64 {
	return r.tag
}

// ErrorCode returns the error code associated with this result.
func (r *Result) ErrorCode() int32 {
	return r.errorCode
}

// ErrorDesc returns the error description associated with this result.
func (r *Result) ErrorDesc() string {
	return r.errorDesc
}

// Value returns the value associated with this result.
func (r *Result) Value() any {
	return r.value
}

// Succeeded returns true if this result represents a successful operation.
func (r *Result) Succeeded() bool {
	return r.errorCode == 0
}

// String returns a string representation of this result.
func (r *Result) String() string {
	return fmt.Sprintf("Result{Key: %s, Tag: %d, ErrorCode: %d, ErrorDesc: "+
		"%s, Value: %v}", r.key, r.tag, r.errorCode, r.errorDesc, r.value)
}
