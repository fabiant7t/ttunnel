package ttunnel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// UnmarshalFromFile attempts to unmarshal the json file into the
// given object.
func UnmarshalFromFile(path string, v interface{}) (err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, &v)
	return
}

// MarshalToFile attempts to write out the given object as json into
// the given file. This function won't overwrite an existing file.
func MarshalToFile(path string, v interface{}) (err error) {
	if fileExists(path) {
		err = fmt.Errorf("Won't overwrite file: %v", path)
		return
	}

	buf, err := json.Marshal(v)
	if err != nil {
		return
	}

	var out bytes.Buffer
	if err = json.Indent(&out, buf, "", "\t"); err != nil {
		return
	}

	err = ioutil.WriteFile(path, out.Bytes(), 0600)
	return
}
