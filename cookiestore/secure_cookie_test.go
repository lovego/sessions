package cookiestore

import (
	"crypto/sha256"
	"fmt"
)

func testSecureCookie(hashKey []byte, name string, value []byte) {
	s := NewSecureCookie(hashKey)
	s.timestampForTest = 1
	encoded, err := s.Encode(name, value)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", encoded)
	plain, err := base64decode(encoded)
	fmt.Printf("%s%v\n", plain[:len(plain)-sha256.Size], plain[len(plain)-sha256.Size:])

	value, err = s.Decode(name, encoded, 7)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", value)
	}

	s.timestampForTest = 8
	_, err = s.Decode(name, encoded[:len(encoded)-1], 7)
	fmt.Println(err)

	s.timestampForTest = 9
	_, err = s.Decode(name, encoded, 7)
	fmt.Println(err)

	_, err = s.Decode(name+"2", encoded, 0)
	fmt.Println(err)

	encoded[0] = '2'
	_, err = s.Decode(name, encoded, 0)
	fmt.Println(err)
}

func ExampleSecureCookie_noHashKey() {
	testSecureCookie(nil, "name", []byte("value"))
	// Output:
	// MXx2YWx1ZXwr5qm7mYaggd_4i5UIfT3M8MF1wpRnYxhR-Q0Sz2W5eA
	// 1|value|[43 230 169 187 153 134 160 129 223 248 139 149 8 125 61 204 240 193 117 194 148 103 99 24 81 249 13 18 207 101 185 120]
	// value
	// illegal base64 data at input byte 52
	// the session has expired
	// the sign is wrong
	// the sign is wrong
}

func ExampleSecureCookie_hasHashKey() {
	testSecureCookie(
		[]byte("test-hash-key"), "name", []byte(`{"userId":1002,"userName":"韩梅梅"}`),
	)

	// Output:
	// MXx7InVzZXJJZCI6MTAwMiwidXNlck5hbWUiOiLpn6nmooXmooUifXw8Gskt8KnwT8hJ31mCyxsbBIqWIET52AV5PhPK5-U9_w
	// 1|{"userId":1002,"userName":"韩梅梅"}|[60 26 201 45 240 169 240 79 200 73 223 89 130 203 27 27 4 138 150 32 68 249 216 5 121 62 19 202 231 229 61 255]
	// {"userId":1002,"userName":"韩梅梅"}
	// illegal base64 data at input byte 96
	// the session has expired
	// the sign is wrong
	// the sign is wrong
}

func ExampleSecureCookie_Encode_tooLong() {
	s := NewSecureCookie(nil)
	s.MaxLength = 30
	_, err := s.Encode("", []byte("value"))
	fmt.Println(err)
	// Output: the encoded value is too long
}

func ExampleSecureCookie_Decode_tooLong() {
	s := NewSecureCookie(nil)
	s.MaxLength = 30
	_, err := s.Decode("", base64encode([]byte("1|v|1234567890123456789012345678901")), 0)
	fmt.Println(err)
	// Output: the value to decode is too long
}

func ExampleSecureCookie_verifyAndRemoveSign() {
	s := NewSecureCookie(nil)

	_, err := s.verifyAndRemoveSign("", []byte("1|v|1234567890123456789012345678901"))
	fmt.Println(err)

	_, err = s.verifyAndRemoveSign("", []byte("|v|12345678901234567890123456789012"))
	fmt.Println(err)

	_, err = s.verifyAndRemoveSign("", []byte("1||12345678901234567890123456789012"))
	fmt.Println(err)

	_, err = s.verifyAndRemoveSign("", []byte("1|v|12345678901234567890123456789012"))
	fmt.Println(err)
	// Output:
	// the value to decode is illegal
	// the value to decode is illegal
	// the value to decode is illegal
	// the sign is wrong
}

func ExampleSecureCookie_verifyAndRemoveTimestamp() {
	s := NewSecureCookie(nil)

	_, err := s.verifyAndRemoveTimestamp([]byte("|v"), 0)
	fmt.Println(err)

	_, err = s.verifyAndRemoveTimestamp([]byte("a|"), 0)
	fmt.Println(err)

	s.timestampForTest = 9
	_, err = s.verifyAndRemoveTimestamp([]byte("1|v"), 7)
	fmt.Println(err)
	// Output:
	// the value to decode is illegal
	// the value to decode is illegal
	// the session has expired
}

func Example_base64encode() {
	fmt.Printf("%s\n", base64encode([]byte("a")))
	fmt.Printf("%s\n", base64encode([]byte("ab")))
	fmt.Printf("%s\n", base64encode([]byte("abc")))
	// Output:
	// YQ
	// YWI
	// YWJj
}
