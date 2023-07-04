package utils

import (
	"encoding/base64"
	"regexp"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid, err := NewUUID()
	if err != nil {
		t.Error("Error encountered during NewUUID() call:", err)
	}

	if len(uuid) != 16 {
		t.Error("Error: UUID byte array expected length: 16. Actual:", len(uuid))
	}

	uuidStr, err := NewUUIDString()
	if err != nil {
		t.Error("Error encountered during NewUUIDString() call:", err)
	}

	isMatch, err := regexp.MatchString("[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}", uuidStr)
	if err != nil {
		t.Error("Error encountered in matching UUID regex to NewUUIDString() result:", err)
	}

	if isMatch == false {
		t.Error("Error: NewUUIDString() result did not succesfully match against UUID regular expression (RFC 4122)")
	}

	uuidStr, err = NewUUIDBase64String()
	if err != nil {
		t.Error("Error encountered during NewUUIDBase64String() call:", err)
	}

	byteArray, err := base64.RawStdEncoding.DecodeString(uuidStr)
	if err != nil {
		t.Error("Error encountered during DecodeString() of a NewUUIDBase64String:", err)
	}

	if len(byteArray) != 16 {
		t.Error("Error: UUID byte array expected length: 16. Actual:", len(uuid))
	}
}
