package keys

import (
	"errors"
	"regexp"
	"strings"
)

const (
	KeySeparator        = ":"
	PrefixSeparator     = "$"
	validatePathPattern = `^([0-9a-z]+[-:]?[a-z0-9.]+){2,}$`

	// Keys 至少必须包含两块内容并以`:`分割.
	// 必须为小写字母+数字
	// 每块内容的字符串长度必须超过1
	// examples:
	// 1. users:profile.10000  2. users:profiles:buddy.10000
)

var (
	validatePathExp *regexp.Regexp

	ErrNotEnoughParts    = errors.New("ErrNotEnoughParts")
	ErrPartStringTooLess = errors.New("ErrPartStringTooLess")
	ErrInvalidKeyFormat  = errors.New("ErrInvalidKeyFormat")
)

func init() {
	validatePathExp = regexp.MustCompile(validatePathPattern)
}

// A Key is a structured abstraction of a redis db identifier.
type Key struct {
	value string
}

// NewKey constructs a new keys from the given fully qualified string.  It will panic if the string isn't valid.  Use
// NewKeyFromString if you don't want to panic.
func NewKey(value string) Key {
	return MustSucceed(NewKeyFromString(value))
}

func NewKeys(values ...string) []Key {
	keys := make([]Key, len(values))
	for k, v := range values {
		keys[k] = NewKey(v)
	}
	return keys
}

// NewKeyFromParts creates a new keys within the current core.Namespace.
func NewKeyFromParts(parts ...string) (key Key, err error) {
	if np := len(parts); np <= 1 {
		err = ErrNotEnoughParts
	} else {
		n := np - 1

		for _, p := range parts {
			l := len(p)
			if l <= 1 {
				err = ErrPartStringTooLess
				return
			}
			n += len(p)
		}

		r := make([]byte, n)
		i := 0

		for _, s := range parts {
			i += copy(r[i:], s)
			i += copy(r[i:], KeySeparator)

		}

		key, err = NewKeyFromString(string(r))
	}

	return
}

// NewKeyFromString creates a new keys within the current core.Namespace from a string
// representation, validating each part.
func NewKeyFromString(value string) (key Key, err error) {
	if namespaceKeyPrefix != "" && !strings.HasPrefix(value, namespaceKeyPrefix) {
		if strings.HasPrefix(value, KeySeparator) {
			value = value[1:]
		}
		value = namespaceKeyPrefix + value
	}

	if !validatePathExp.MatchString(value) {
		err = ErrInvalidKeyFormat
	} else {
		key.value = value
	}
	return
}

// NewKeyFromBytesUnchecked creates a new keys from the given byte array without validation of any kind.
func NewKeyFromBytesUnchecked(value []byte) Key {
	return Key{string(value)}
}

// NewKeyFromStringUnchecked creates a new keys from the given byte array without validation of any kind.
func NewKeyFromStringUnchecked(value string) Key {
	return Key{value}
}

// MustSucceed is a handy wrapper for panicking when creating a new Key.  Example: MustSucceed(NewKeyFromString("xyz")).
func MustSucceed(key Key, err error) Key {
	if err != nil {
		panic(err)
	} else {
		return key
	}
}

// Clear any value in this keys.
func (k *Key) Clear() {
	k.value = ""
}

// Parts returns a slice of the component parts of the keys.
func (k Key) Parts() []string {
	return strings.Split(k.String(), KeySeparator)
}

//// Base returns the final KeySeparator-separated element of the keys.
//func (k Key) Base() string {
//	return k.String()[strings.LastIndex(k.String(), KeySeparator)+1:]
//}

//// Prefixes returns a slice of PrefixSeparator-separated elements of the Base.
//func (k Key) Prefixes() []string {
//	p := strings.Split(k.Base(), PrefixSeparator)
//	if len(p) > 0 {
//		p = p[:len(p)-1]
//	}
//	return p
//}

func (k Key) String() string {
	return k.value
}

func (k Key) Bytes() []byte {
	return []byte(k.value)
}

func (k Key) IsEqual(s string) bool {
	return k.String() == s
}

func (k Key) IsEmpty() bool {
	return k.value == ""
}

func (k Key) HasValue() bool {
	return k.value != ""
}

type ChangesType string

const (
	Set     ChangesType = "set"
	HSet    ChangesType = "hset"
	Delete  ChangesType = "del"
	Expired ChangesType = "expired"
)
