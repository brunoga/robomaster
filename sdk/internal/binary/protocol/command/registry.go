package command

import "fmt"

var (
	keyToCommandRequestMap  map[uint16]Command = make(map[uint16]Command)
	keyToCommandResponseMap map[uint16]Command = make(map[uint16]Command)
)

// Register registers a command with the given set and id. The command must
// implement either the Request or Response interface or it panics.
func Register(cmdSet, cmdID byte, cmd Command) {
	switch r := cmd.(type) {
	case Request:
		register(cmdSet, cmdID, true, r)
	case Response:
		register(cmdSet, cmdID, false, r)
	default:
		panic(fmt.Sprintf("unexpected command type: %T", r))
	}
}

// GetRequest returns a new request command with the given set and id if one is
// found. Otherwise it panics.
func GetRequest(cmdSet, cmdID byte, data []byte) Request {
	return get(cmdSet, cmdID, data, true).(Request)
}

// GetResponse returns a new response command with the given set and id if one
// is found. Otherwise it panics.
func GetResponse(cmdSet, cmdID byte, data []byte) Response {
	return get(cmdSet, cmdID, data, false).(Response)
}

func getMapAndName(isRequest bool) (map[uint16]Command, string) {
	var keyToCommandMap map[uint16]Command
	var name string
	if isRequest {
		keyToCommandMap = keyToCommandRequestMap
		name = "request"
	} else {
		keyToCommandMap = keyToCommandResponseMap
		name = "response"
	}

	return keyToCommandMap, name
}

func register(cmdSet, cmdID byte, isRequest bool, cmd Command) {
	keyToCommandMap, name := getMapAndName(isRequest)

	key := Key(cmdSet, cmdID)

	if _, ok := keyToCommandMap[key]; ok {
		panic(fmt.Sprintf(
			"command %s already registered: cmdSet=%02x, cmdID=%02x",
			name, cmdSet, cmdID))
	}

	keyToCommandMap[key] = cmd
}

func get(cmdSet, cmdID byte, data []byte, isRequest bool) Command {
	keyToCommandMap, name := getMapAndName(isRequest)

	key := Key(cmdSet, cmdID)

	if cmd, ok := keyToCommandMap[key]; ok {
		return cmd.New(data)
	}

	panic(fmt.Sprintf("command %s not registered: cmdSet=%x, cmdID=%x",
		name, cmdSet, cmdID))
}
