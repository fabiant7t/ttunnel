package ttunnel

import (
	"bytes"
	"os"
	"testing"
)

func TestClientConfig(t *testing.T) {
	name := "testClientConfig"

	defer os.Remove(ClientConfigPath(name))

	cc := ClientConfig{}
	cc.Host = "localhost:1044"
	cc.Token = []byte("A token")
	cc.Port = 1024

	// Write the config.
	if err := WriteClientConfig(name, cc); err != nil {
		t.Fatal(err)
	}

	// Should refuse overwrite.
	if err := WriteClientConfig(name, cc); err == nil {
		t.Fatal("Should have refused overwrite.")
	}

	// Read the config.
	cc2, err := ReadClientConfig(name)
	if err != nil {
		t.Fatal(err)
	}

	// Check configs are the same.
	if cc.Host != cc2.Host || cc.Port != cc2.Port {
		t.Fatal("Configs don't match.")
	}
	if !bytes.Equal(cc.Token, cc2.Token) {
		t.Fatal("Tokens don't match.")
	}
}
