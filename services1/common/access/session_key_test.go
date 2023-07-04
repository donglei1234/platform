package access

import (
	"testing"

	"github.com/donglei1234/platform/services/common/nosql/document"
)

func TestSessionKey(t *testing.T) {
	document.SetNamespace("test")

	// test 1 - generate a generic key
	skAccess := AccessLevel(1)
	skMeta := 0x2
	skAuth := SessionAuthority(3)

	sk := NewSessionKey(skAccess, skMeta, skAuth)
	t.Logf("new key - %s\n", sk.String())

	// test 1b - ensure that sk HAS a key set
	if sk.Key.IsEmpty() {
		t.Fatalf("Unexpected empty nosql.Key in SessionKey %v\n", sk)
	}

	// test 1c - ensure that metadata fileds are as expected
	if sk.AccessLevel != skAccess {
		t.Fatalf("Got unexpected access level %v, expected %v\n", sk.AccessLevel, skAccess)
	}
	if sk.Metadata != skMeta {
		t.Fatalf("Got unexpected metadata %v, expected %v\n", sk.Metadata, skMeta)
	}
	if sk.Authority != skAuth {
		t.Fatalf("Got unexpected authority %v, expected %v\n", sk.Authority, skAuth)
	}

	// test 2 - parse the generated token
	sk2 := ParseSessionKeyFromToken(sk.String())
	t.Logf("parsed - %s\n", sk2.String())

	// test 2b - ensure that sk2 does NOT have a key set
	if !sk2.Key.IsEmpty() {
		t.Errorf("Unexpected nosql.Key %v in SessionKey %v\n", sk2.Key, sk2)
	}

	// test 2c - ensure that sk1 and sk2 look the same
	if sk.String() != sk2.String() {
		t.Errorf("Parsed key %v does not match original key %v\n", sk2, sk)
	}

	// test 3 - legacy parsing
	var sk3 SessionKey
	tok := generateSessionToken()
	if k, e := document.NewKeyFromParts(tok); e != nil {
		t.Errorf(e.Error())
	} else {
		sk3 = ParseSessionKey(k)
	}
	t.Logf("legacy - %s\n", sk3.String())

	// test 3b - verify legacy flag
	if !sk3.IsLegacy() {
		t.Errorf("SessionKey %v is not flagged as legacy\n", sk3)
	}

	// test 3c - verify zeroed metadata flags
	if sk3.AccessLevel != 0 || sk3.Metadata != 0 || sk3.Authority != Unknown_SessionAuthority {
		t.Errorf("SessionKey %v has non-zero metadata fields\n", sk3)
	}
}
