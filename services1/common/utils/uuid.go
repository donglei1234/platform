package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// generates a random UUID according to RFC 4122
func NewUUID() ([16]byte, error) {
	uuid := [16]byte{}

	n, err := io.ReadFull(rand.Reader, uuid[:])
	if n != len(uuid) || err != nil {
		return uuid, err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuid, nil
}

func NewUUIDString() (string, error) {
	if uuid, err := NewUUID(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
	}
}

func NewUUIDBase64String() (string, error) {
	if uuid, err := NewUUID(); err != nil {
		return "", err
	} else {
		return base64.RawStdEncoding.EncodeToString(uuid[:]), nil
	}
}
