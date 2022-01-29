package dsp

import (
	"io/ioutil"
)

// UnwrapData unwraps the data in the DSP file pointed by filename. Returns the
// data and a nil error on success or nil and an error on failure.
func UnwrapData(fileName string) ([]byte, error) {
	return decodeDsp(fileName)
}

// WrapData wraps the data in the file pointed by filename in a DSP container.
// Returns the container wrapped data and a nil error on success or nil and an
// error on failure.
func WrapData(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return encodeDsp(data)
}
