package cookiestore

import (
	"fmt"
	"net/http"
)

func ExampleCookieStore_Get() {
	s := New("test-hash-key")
	type session struct {
		UserId   int    `json:"userId"`
		UserName string `json:"userName"`
	}
	var data session

	ck := &http.Cookie{
		Name:  "name",
		Value: "MXx7InVzZXJJZCI6MTAwMiwidXNlck5hbWUiOiLpn6nmooXmooUifXw8Gskt8KnwT8hJ31mCyxsbBIqWIET52AV5PhPK5-U9_w",
	}
	data = session{}
	if err := s.Get(ck, &data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}

	ck = &http.Cookie{
		Name:  "name",
		Value: "xyz",
	}
	data = session{}
	if err := s.Get(ck, &data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}

	ck = &http.Cookie{
		Name:  "name",
		Value: "",
	}
	data = session{}
	if err := s.Get(ck, &data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}

	// Output:
	// {1002 韩梅梅}
	// the value to decode is illegal
	// {0 }
}

func ExampleCookieStore_Save() {
	ck := &http.Cookie{Name: "name", MaxAge: 86400}
	data := map[string]interface{}{"userId": 1002, "userName": "韩梅梅"}

	rw := testResponseWriter{make(http.Header)}
	s := New("test-hash-key")
	s.SecureCookie.timestampForTest = 1
	if err := s.Save(rw, ck, data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rw.header)
	}

	s.SecureCookie.MaxLength = 50
	fmt.Println(s.Save(rw, ck, data))

	// Output:
	// map[Set-Cookie:[name=MXx7InVzZXJJZCI6MTAwMiwidXNlck5hbWUiOiLpn6nmooXmooUifXw8Gskt8KnwT8hJ31mCyxsbBIqWIET52AV5PhPK5-U9_w; Max-Age=86400]]
	// the encoded value is too long
}

func ExampleCookieStore_Delete() {
	ck := &http.Cookie{Name: "name", MaxAge: 1}

	rw := testResponseWriter{make(http.Header)}
	s := New("")
	s.Delete(rw, ck)
	fmt.Println(rw.header)
	// Output:
	// map[Set-Cookie:[name=; Expires=Thu, 01 Jan 1970 00:00:01 GMT; Max-Age=0]]
}

func ExampleCookieStore_EncodeData() {
	s := New("")
	fmt.Println(s.EncodeData("name", func() {}))

	// Output:
	// [] json: unsupported type: func()
}

type testResponseWriter struct {
	header http.Header
}

func (rw testResponseWriter) Header() http.Header {
	return rw.header
}

func (rw testResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (rw testResponseWriter) WriteHeader(int) {
}
