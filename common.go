package ttunnel

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var ConfigDir string  // Configuration directory.
var ConfigPath string // The server configuration file path.
var TunnelDir string  // Directory containing tunnels.
var RemovedDir string // Removed client directory.
var CertPath string   // Path to X.509 certificate file.
var KeyPath string    // Path to X.509 key file.

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Initialize directory paths
	ConfigDir = filepath.Join(usr.HomeDir, ".ttunnel")
	ConfigPath = filepath.Join(ConfigDir, "server.json")
	TunnelDir = filepath.Join(ConfigDir, "tunnels")
	RemovedDir = filepath.Join(ConfigDir, "removed")

	// Create directories.
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(TunnelDir, 0700); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(RemovedDir, 0700); err != nil {
		panic(err)
	}

	// Initialize file paths.
	CertPath = filepath.Join(ConfigDir, "server.crt")
	KeyPath = filepath.Join(ConfigDir, "server.key")
}

// fileExists checks if a path exists. Arguments are joined with
// filepath.Join to construct the full path.
func FileExists(elem ...string) bool {
	path := filepath.Join(elem...)
	_, err := os.Stat(path)
	return err == nil
}

// RandomBytes return cryptographically secure random bytes.
func RandomBytes(length int) (key []byte, err error) {
	key = make([]byte, length)
	_, err = io.ReadFull(rand.Reader, key)
	return
}

// Helper function to unmarshal json from a file path.
func UnmarshalFrom(path string, v interface{}) (err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &v)

	return
}

// Helper function to marshal json to a file path.
func MarshalTo(path string, v interface{}) (err error) {
	if FileExists(path) {
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
