package sessions

import (
	"fmt"
	"net/http"
)

func ExampleCookieStore_Get() {
	cs := NewCookieStore("test-hash-key")
	req := http.Request{
		Header: http.Header{"Cookie": []string{
			"name=MXx7InVzZXJJZCI6MTAwMiwidXNlck5hbWUiOiLpn6nmooXmooUifXw8Gskt8KnwT8hJ31mCyxsbBIqWIET52AV5PhPK5-U9_w",
		}},
	}
	sess, err := cs.Get(&req, "name", 0)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", sess.Cookie)
		fmt.Printf("%s\n", sess.data)
	}

	req = http.Request{Header: http.Header{"Cookie": []string{"name=xyz"}}}
	fmt.Println(cs.Get(&req, "name", 0))
	fmt.Println(cs.Get(&req, "none", 0))

	// Output:
	// name=MXx7InVzZXJJZCI6MTAwMiwidXNlck5hbWUiOiLpn6nmooXmooUifXw8Gskt8KnwT8hJ31mCyxsbBIqWIET52AV5PhPK5-U9_w
	// {"userId":1002,"userName":"韩梅梅"}
	// <nil> the value to decode is illegal
	// <nil> <nil>
}

func ExampleCookieStore_Save() {
	ck := &http.Cookie{Name: "name", MaxAge: 86400}
	sess := &Session{Cookie: ck, Data: map[string]interface{}{"userId": 1002, "userName": "韩梅梅"}}

	rw := testResponseWriter{make(http.Header)}
	cs := NewCookieStore("test-hash-key")
	cs.SecureCookie.timestampForTest = 1
	if err := cs.Save(rw, sess); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rw.header)
	}

	cs.SecureCookie.MaxLength = 50
	fmt.Println(cs.Save(rw, sess))

	// Output:
	// map[Set-Cookie:[name=MXx7InVzZXJJZCI6MTAwMiwidXNlck5hbWUiOiLpn6nmooXmooUifXw8Gskt8KnwT8hJ31mCyxsbBIqWIET52AV5PhPK5-U9_w; Max-Age=86400]]
	// the encoded value is too long
}

func ExampleCookieStore_Delete() {
	ck := &http.Cookie{Name: "name", MaxAge: 1}
	sess := &Session{Cookie: ck}

	rw := testResponseWriter{make(http.Header)}
	cs := NewCookieStore("")
	cs.Delete(rw, sess)
	fmt.Println(rw.header)
	// Output:
	// map[Set-Cookie:[name=; Expires=Thu, 01 Jan 1970 00:00:01 GMT; Max-Age=0]]
}

func ExampleCookieStore_DecodeData() {
	ck := &http.Cookie{Name: "name", MaxAge: 1}
	sess := &Session{Cookie: ck, Data: func() {}}
	cs := NewCookieStore("")
	fmt.Println(cs.EncodeData(sess))

	// Output: json: unsupported type: func()
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
