package result

import (
	"reflect"
	"testing"

	"github.com/brunoga/robomaster/sdk2/unitybridge/unity/key"
	"github.com/brunoga/robomaster/sdk2/unitybridge/unity/result/value"
)

func TestNewFromJSON(t *testing.T) {
	type args struct {
		jsonData []byte
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		{
			name: "empty json data",
			args: args{
				jsonData: []byte{},
			},
			want: &Result{
				errorCode: -1,
				errorDesc: "empty or nil json data",
			},
		},
		{
			name: "nil json data",
			args: args{
				jsonData: nil,
			},
			want: &Result{
				errorCode: -1,
				errorDesc: "empty or nil json data",
			},
		},
		{
			name: "error unmarshalling json data",
			args: args{
				jsonData: []byte("invalid"),
			},
			want: &Result{
				errorCode: -1,
				errorDesc: "error unmarshalling json data: invalid character " +
					"'i' looking for beginning of value",
			},
		},
		{
			name: "valid json data",
			args: args{
				jsonData: []byte(`{"key":117440513,"tag":0,"value":{"value":true}}`),
			},
			want: &Result{
				key:   key.KeyAirLinkConnection,
				value: &value.Bool{Value: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromJSON(tt.args.jsonData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
