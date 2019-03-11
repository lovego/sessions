package cookiestore

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"time"
)

var ErrorEncodedValueTooLong = errors.New("the encoded value is too long")
var ErrorValueToDecodeToolong = errors.New("the value to decode is too long")
var ErrorValueToDecodeIllegal = errors.New("the value to decode is illegal")
var ErrorSignWrong = errors.New("the sign is wrong")
var ErrorSessionExpired = errors.New("the session has expired")

type SecureCookie struct {
	hashKey          []byte
	MaxLength        int
	timestampForTest int64
}

func NewSecureCookie(hashKey []byte) *SecureCookie {
	s := &SecureCookie{
		hashKey:   hashKey,
		MaxLength: 4096,
	}
	return s
}

// Encode encodes a cookie value.
//
// It signs the value with a message authentication code, and encodes it.
//
// The name argument is the cookie name. It is signed with the encoded value.
// The value argument is the value to be encoded.
func (s *SecureCookie) Encode(name string, value []byte) ([]byte, error) {
	b := []byte(fmt.Sprintf("%s|%d|%s|", name, s.nowTimestamp(), value))
	b = append(b, s.createSign(b)...) // append sign
	b = b[len(name)+1:]               // remove name
	b = base64encode(b)
	if s.MaxLength != 0 && len(b) > s.MaxLength {
		return nil, ErrorEncodedValueTooLong
	}
	return b, nil
}

// Decode decodes a cookie value.
//
// It decodes the value, and verifies a message authentication code.
//
// The name argument is the cookie name. It must be the same name used when it was stored.
// The value argument is the encoded cookie value.
// The maxAge argument is the max seconds since the value is generated.
func (s *SecureCookie) Decode(name string, value []byte, maxAge int64) ([]byte, error) {
	if s.MaxLength != 0 && len(value) > s.MaxLength {
		return nil, ErrorValueToDecodeToolong
	}
	b, err := base64decode(value)
	if err != nil {
		return nil, err
	}
	// now b is "timestamp|value|sign".
	b, err = s.verifyAndRemoveSign(name, b)
	if err != nil {
		return nil, err
	}
	// now b is "timestamp|value".
	return s.verifyAndRemoveTimestamp(b, maxAge)
}

func (s *SecureCookie) createSign(b []byte) []byte {
	var h hash.Hash
	if len(s.hashKey) > 0 {
		h = hmac.New(sha256.New, s.hashKey)
	} else {
		h = sha256.New()
	}
	h.Write(b)
	return h.Sum(nil)
}

func (s *SecureCookie) verifyAndRemoveSign(name string, b []byte) ([]byte, error) {
	// now b is "timestamp|value|sign".
	pos := len(b) - sha256.Size // position of the first byte of the signature.
	// there must be at lease 4 byte before pos, example:  "0|v|"
	if pos < 4 || b[pos-1] != '|' {
		return nil, ErrorValueToDecodeIllegal
	}
	gotSign := b[pos:]
	b = b[:pos]
	// now b is "timestamp|value|".
	realSign := s.createSign(append([]byte(name+"|"), b...))
	if subtle.ConstantTimeCompare(gotSign, realSign) != 1 {
		return nil, ErrorSignWrong
	}
	return b[:len(b)-1], nil
}

func (s *SecureCookie) verifyAndRemoveTimestamp(b []byte, maxAge int64) ([]byte, error) {
	// now b is "timestamp|value".
	pos := bytes.IndexByte(b, '|')
	if pos < 1 {
		return nil, ErrorValueToDecodeIllegal
	}
	gotTimestamp, err := strconv.ParseInt(string(b[:pos]), 10, 64)
	if err != nil {
		return nil, ErrorValueToDecodeIllegal
	}
	if maxAge > 0 && s.nowTimestamp()-gotTimestamp > maxAge {
		return nil, ErrorSessionExpired
	}
	return b[pos+1:], nil
}

func (s *SecureCookie) nowTimestamp() int64 {
	if s.timestampForTest > 0 {
		return s.timestampForTest
	}
	return time.Now().UTC().Unix()
}

func base64encode(value []byte) []byte {
	encoded := make([]byte, base64.RawURLEncoding.EncodedLen(len(value)))
	base64.RawURLEncoding.Encode(encoded, value)
	return encoded
}

func base64decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.RawURLEncoding.DecodedLen(len(value)))
	n, err := base64.RawURLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:n], nil
}
