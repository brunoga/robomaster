package dji

import (
	"reflect"
	"testing"
)

func TestNewResultFromJSON_NilData(t *testing.T) {
	result := NewResultFromJSON(nil)

	expected := &Result{
		0,
		0,
		nil,
		-1,
		"empty or nil json data",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected %#+v, got %#+v", expected, result)
	}

}

func TestNewResultFromJSON_EmptyData(t *testing.T) {
	result := NewResultFromJSON([]byte(""))

	expected := &Result{
		0,
		0,
		nil,
		-1,
		"empty or nil json data",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected %#+v, got %#+v", expected, result)
	}

}

func TestNewResultFromJSON_InvalidData(t *testing.T) {
	result := NewResultFromJSON([]byte("invalid"))

	expected := &Result{
		0,
		0,
		nil,
		-1,
		"invalid json data: invalid",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected %#+v, got %#+v", expected, result)
	}

}


func TestNewResultFromJSON_Error(t *testing.T) {
	jsonData := []byte("{\"Value\":{\"value\":true}, \"Error\":-1, \"Key\":117440513, \"Tag\":0}")
	result := NewResultFromJSON(jsonData)

	expected := &Result{
		1,
		0,
		nil,
		-1,
		"result error",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected %#+v, got %#+v", expected, result)
	}
}

func TestNewResultFromJSON_Success(t *testing.T) {
	jsonData := []byte("{\"Value\":{\"value\":true}, \"Error\":0, \"Key\":117440513, \"Tag\":0}")
	result := NewResultFromJSON(jsonData)

	expected := &Result{
		1,
		0,
		true,
		0,
		"",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected %#+v, got %#+v", expected, result)
	}
}