package ttunnel

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"fmt"
	"time"
)

// A Token is encrypted and signed by the server, and is used by a
// client to initiate a forwarding connection.
type Token struct {
	Name        string // The client's name. Used for logging and blocking.
	ConnectAddr string // The address to connect to for the client.
	Expires     int64  // The unix timestamp for expiration.
}

// A TokenEncoder is used to encode and decode tokens.
type TokenHandler struct {
	aead cipher.AEAD
}

// NewTokenEncoder creates a new encoder using the given key. The key
// must be 32 bytes in length.
func NewTokenHandler(key []byte) (th *TokenHandler, err error) {
	th = new(TokenHandler)

	// Check key lengths.
	if len(key) != 32 {
		err = fmt.Errorf("Key must be 32 bytes, not %v.", len(key))
		return
	}

	// Initialize an AES cipher.
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// Create aead.
	th.aead, err = cipher.NewGCM(block)
	if err != nil {
		return
	}

	return
}

// Encode takes a token and returns an encoded byte array.
func (th *TokenHandler) Encode(t Token) (encToken []byte, err error) {
	buf := new(bytes.Buffer)

	// Encode value into a byte string.
	if err = gob.NewEncoder(buf).Encode(t); err != nil {
		return
	}

	plaintext := buf.Bytes()
	nonce, err := RandomBytes(th.aead.NonceSize())
	if err != nil {
		return
	}

	// Encrypt the plaintext.
	encToken = th.aead.Seal(nil, nonce, plaintext, nil)
	encToken = append(nonce, encToken...)

	return
}

// Decode takes an encoded byte array and returns a Token.  An error
// will be returned if the token is expired or the name is blocked.
func (th *TokenHandler) Decode(encToken []byte) (token Token, err error) {

	if len(encToken) < th.aead.NonceSize() {
		err = fmt.Errorf("Token is too short: %v", len(encToken))
		return
	}

	nonce := encToken[:th.aead.NonceSize()]
	encToken = encToken[th.aead.NonceSize():]

	// Decrypt ciphertext.
	plaintext, err := th.aead.Open(nil, nonce, encToken, nil)
	if err != nil {
		return
	}

	// Decode plaintext into object.
	err = gob.NewDecoder(bytes.NewBuffer(plaintext)).Decode(&token)
	if err != nil {
		return
	}

	return
}

// Verify checks that a token is not expired and that it's allowed by
// the server.
func (th *TokenHandler) Verify(token Token) (err error) {
	if token.Expires <= time.Now().Unix() { // Expired token.
		err = fmt.Errorf("Token expired %v days ago.",
			(time.Now().Unix()-token.Expires)/(3600*24))
	} else if !FileExists(ClientConfigPath(token.Name)) { // Client blocked.
		err = fmt.Errorf("Client not allowed: %v", token.Name)
	}
	return
}
