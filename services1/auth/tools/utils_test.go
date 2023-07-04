package tools

import (
	"testing"
)

const Jwt_Key = `MIIEpAIBAAKCAQEArsRdY906rO9mJ8kz9qnK7LebHNsLDu7mmv+kmJbNMRWE4Lg4
dzEz1w+kEDkTuIPJb7omwifLIfmfXA55Y4J0vGCVsUYJS6e5n0XJBRUqZyKJfqWW
0IHYjFKBVUMv2uDSz6qxuglgaGmthpm+K1l9YlUugNchXO8R9oRtw14xPQy0mLyP
+z1ZG6ubfiYWAwgesXphrCAkwfSvq30dCxJ2ssD30lfe6k63T4/XZZ8iKWy2dknG
x/i4im3UwaFHMIE5KfVDRnfXIRBjOfGoBM2pz/j4QvcBNbCekR+F5C6XmFsWfnJF
+OOrleXs+Uytr8Ie4qqdBAkokjhIXh21gKzY+wIDAQABAoIBAH1VpgQwbBwJtCFk
UjfbnQQWyM7w6AVVn1wZallkDNPesuTWOOiCCMN7HBXmWRZCrPvLbHBhSXScKIVE
fBm5PS67tZ7ks1Xme2CWE4vpmfdM9X42EuqFtF8t7ELRZEh8Y7M6nyrj/pfi8edN
6uv7ycfzft43al3TtfYtEurBal50xrqn2MRgNEzClM95kqQjEiolH44JkLuaS2ZG
gl/25kXizN7LU/7fawCahY3eOcaajgpQgdY08T7Ac7e3phkc5wi01vYdIqFWANTr
K1Bm17LYD+YSEliEGtWhBhS6zC/0iMRJmkdSrAfjHOMO/ZxxB/FBwH820Lj+4BkH
2112ecECgYEA2BgGJR0Q0PaNICVqrJYUg89YoD0mfwSs6u3QWqEwOSyrnvVjKIzQ
s9zy2liyp1iuBZVbC7JBYZ0vIW5EPtpBv67QHiZ2wFvq24aa4qXJScF7WjDJC/nu
FmWBT+cDPk7Fq8Dhq7/V8G9UTCNrPQ6EU0JGzO2ZXrKR9yx0JVdTNZ0CgYEAzwqX
BnME4BhFTFGNFKRkCVDyE2xl78muvYxLaRghkEoILDDKx02+Ro1hpmwdBg2YRsq3
rHNobd7QAgUsJ2a/D02Ws69BOfyR41JdLNdQf9AuNbhnx2wt9+LIuCh6QdRSIlCz
8vqgxtCOufG/RiBEsys1e88yN9tLXw0ahxzVkXcCgYEAgBKtAa4aW17lAZprgOJq
QjzPsBjOChYBTjOoey5xYFGDXfDd1hivUQqwvIw5RkmeyhxdG6+IZIw+dLffpsjA
kxTOsC/nzdYsefNrNM1BYX9U1n13iquUWp3KaErodzNEoKjul1ZZO+kfswiC9Gr4
LkhnoeloLuHy9OXZ6I+691kCgYAD2D1/uCkSBEFdLnKBPKBC8Abex3eJIWSmTnGk
DPeql0VZlLpfQxlSdBOpJH7RevUl82O/xxDcYXPZZcExJh5MKXOv+IQskH0hDImA
aKstBG+nPbpnoKGl5cc2bEIA8PFpg7gjvfW4W20bFNspcTX1YcsHdoyHLwJSTIr+
kVGB9wKBgQDRXv4wxGqiySBzzosjvxGV0RneLMzfJPU6dvSyJ3dMPFiFnRJ+RqPR
VR5Zbc1Djz1LceQRf/wfm79UuPTUlnfir5viylj08dMhcZ8OTU6rIKZGXS+gvkyR
NIoKwAGWBiJ3tU8PGjht3Sj6n5gchOYBl/3Uvh/eSTKaBqIVHyjKUQ==
-----END RSA PRIVATE KEY-----`

func TestCreatJwt(t *testing.T) {
	type testCase struct {
		profileId     string
		errorExpected error
	}
	// 对封装的uid没有校验格式，所以目前string类型都可被写入
	var testCases = []testCase{
		{
			profileId:     "",
			errorExpected: nil,
		},
		{
			profileId:     "111",
			errorExpected: nil,
		},
		{
			profileId:     "$%^&&*(*))",
			errorExpected: nil,
		},
		{
			profileId:     "------111111111----",
			errorExpected: nil,
		},
	}

	for i, tc := range testCases {
		_, err := CreatProfileJwt(tc.profileId, Jwt_Key)
		if err == tc.errorExpected {
			continue
		} else {
			t.Fatal("Unexpected error encountered in test case #", i+1, ":", err)
		}
	}
}

func TestParseToken(t *testing.T) {
	type testCase struct {
		expectedKey   string
		errorExpected error
	}
	var testCases = []testCase{
		{
			expectedKey:   "",
			errorExpected: ErrTokenMalformed,
		},
		{
			expectedKey:   "111",
			errorExpected: ErrTokenMalformed,
		},
		{
			expectedKey:   "$%^&&*(*))",
			errorExpected: ErrTokenMalformed,
		},
		{
			expectedKey:   "------111111111----",
			errorExpected: ErrTokenMalformed,
		},
		{
			expectedKey:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.lskadjhfkldsgfhdekljghdfkjg.sadjhfgsdhjfgdsjhfg",
			errorExpected: ErrTokenMalformed,
		},
		{
			expectedKey:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA5MTkwMTAsInVpZCI6IjQ0MDY5MDcyNiJ9.ph1w3pRFiFL3xUpNrhcjhyZzKcG2FkbB2aatbN20rtI",
			errorExpected: ErrTokenExpired,
		},
		{
			expectedKey:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA5MzM4MDUsInVpZCI6IjM3ODk2NzA3MjYifQ.9HtLpjlOiQTLeRQ46vOnL0nlnTCJez88fNEFnQu9IY0",
			errorExpected: ErrTokenExpired,
		},
		{
			expectedKey:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDEwMzM0OTksInVpZCI6IiQlXlx1MDAyNlx1MDAyNiooKikpIn0.MDHqTLZUMLpv8NMvG5Hot2Zj9d7p9D68GyXq45cvv9k",
			errorExpected: nil,
		},
		{
			expectedKey:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDEwMzM0OTksInVpZCI6Ii0tLS0tLTExMTExMTExMS0tLS0ifQ.xUlXVJVtGUTxGezCMWSPJ8dnOFo_6GdAz9Q1WtOh6BQ",
			errorExpected: nil,
		},
	}

	for i, tc := range testCases {
		_, err := ParseProfileToken(tc.expectedKey, Jwt_Key)
		if err == tc.errorExpected {
			continue
		} else {
			t.Fatal("Unexpected error encountered in test case #", i+1, ":", err)
		}
	}
}
