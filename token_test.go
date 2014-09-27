package ttunnel

import (
	"os"
	"testing"
	"time"
)

func TestToken1(t *testing.T) {
	name := "TestClient234"

	// Clean up from old tests.
	os.Remove(ClientConfigPath(name))
	os.Remove(ClientRemovedPath(name))

	key, err := RandomBytes(32)
	if err != nil {
		t.Fatal(err)
	}

	// Short key.
	_, err = NewTokenHandler(key[:31])
	if err == nil {
		t.Fatal("Should have failed with short key.")
	}

	// New, ok.
	te, err := NewTokenHandler(key)
	if err != nil {
		t.Fatal(err)
	}

	t1 := Token{}
	t1.Name = name
	t1.ConnectAddr = "localhost:2020"
	t1.Expires = time.Now().Unix() + 30

	// Encode.
	t1Enc, err := te.Encode(t1)
	if err != nil {
		t.Fatal(err)
	}

	// Short ciphertext.
	_, err = te.Decode(t1Enc[:28])
	if err == nil {
		t.Fatal("Should have failed with short ciphertext.")
	}

	// Bad ciphertext.
	t1Enc[0] += 3
	_, err = te.Decode(t1Enc)
	if err == nil {
		t.Fatal("Should have failed with modified ciphertext.")
	}

	t1Enc[0] -= 3

	// Clipped ciphertext.
	_, err = te.Decode(t1Enc[:len(t1Enc)-1])
	if err == nil {
		t.Fatal("Should have failed with clipped ciphertext.")
	}

	// Decode.
	t2, err := te.Decode(t1Enc)
	if err != nil {
		t.Fatal(err)
	}

	if t1.Name != t2.Name ||
		t1.ConnectAddr != t2.ConnectAddr ||
		t1.Expires != t2.Expires {
		t.Fatal("Tokens don't match.")
	}

	// Test verification.
	t1.Expires = time.Now().Unix() + 1

	t1Enc, err = te.Encode(t1)
	if err != nil {
		t.Fatal(err)
	}

	t2, err = te.Decode(t1Enc)
	if err != nil {
		t.Fatal(err)
	}

	// Create clien config so we have a valid token.
	cfg := ClientConfig{}
	cfg.Host = "localhost:1033"
	cfg.Token = t1Enc
	cfg.Port = 4000

	if err = WriteClientConfig(t2.Name, cfg); err != nil {
		t.Fatal(err)
	}

	// Verification should succeed.
	if err = te.Verify(t2); err != nil {
		t.Fatal(err)
	}

	// Delete the client file.
	if err = RemoveClientConfig(t2.Name); err != nil {
		t.Fatal(err)
	}

	// Client file doesn't exist.
	if err = te.Verify(t2); err == nil {
		t.Fatal("Client shouldn't be allowed.")
	}

	// Wait for expiration.
	time.Sleep(1100 * time.Millisecond)

	if err = te.Verify(t2); err == nil {
		t.Fatal("Should have expired.")
	}
}
