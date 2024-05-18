package result

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/brunoga/robomaster/unitybridge/unity/key"
)

// Reesult represents a result from an operation on a key. The zero value
// is not valid.
type Result struct {
	key       *key.Key
	tag       uint64
	errorCode int64
	errorDesc string
	value     any
}

type jsonResult struct {
	Key   uint32
	Tag   uint64
	Error int64
	Value json.RawMessage // defer decoding value until we know the type
}

// New creates a new Result with the given parameters.
func New(key *key.Key, tag uint64, errorCode int64, errorDesc string,
	value any) *Result {
	if reflect.TypeOf(key.ResultValue()) != reflect.TypeOf(value) {
		panic(fmt.Sprintf("result value type (%s) does not match key %s value "+
			"type (%s)", reflect.TypeOf(value), key, reflect.TypeOf(
			key.ResultValue())))
	}

	return &Result{
		key:       key,
		tag:       tag,
		errorCode: errorCode,
		errorDesc: errorDesc,
		value:     value,
	}
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

	err := json.Unmarshal(jsonData, &r)
	if err != nil {
		r.errorCode = -1
		r.errorDesc = fmt.Sprintf("error unmarshalling json data: %s",
			err.Error())
		return r
	}

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
func (r *Result) ErrorCode() int64 {
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

func (r *Result) UnmarshalJSON(data []byte) error {
	jr := jsonResult{}

	err := json.Unmarshal(data, &jr)
	if err != nil {
		return err
	}

	key, err := key.FromSubType(jr.Key)
	if err != nil {
		return err
	}

	// First try to parse value as a string (for when a request has an empty
	// value response)
	var stringValue string
	err = json.Unmarshal(jr.Value, &stringValue)
	if err != nil || len(stringValue) != 0 {
		// Parsing as a string either failed or resulted in an non-empty string.
		// Try to parse as the actual value type.
		value := key.ResultValue()
		err = json.Unmarshal(jr.Value, &value)
		if err != nil {
			return err
		}

		r.value = value
	}

	errorDesc := ""
	if jr.Error != 0 {
		errorDesc = fmt.Sprintf("error %d", jr.Error)
	}

	r.key = key
	r.tag = jr.Tag
	r.errorCode = jr.Error
	r.errorDesc = errorDesc

	return nil
}

func (r *Result) MarshalJSON() ([]byte, error) {
	if reflect.TypeOf(r.Key().ResultValue()) != reflect.TypeOf(r.value) {
		return nil, fmt.Errorf("result value type (%s) does not match key %s "+
			"value type (%s)", reflect.TypeOf(r.value), r.Key(),
			reflect.TypeOf(r.Key().ResultValue()))
	}

	value, err := json.Marshal(r.value)
	if err != nil {
		return nil, err
	}

	jr := jsonResult{
		Key:   r.key.SubType(),
		Tag:   r.tag,
		Error: r.errorCode,
		Value: json.RawMessage(value),
	}

	return json.Marshal(jr)
}
