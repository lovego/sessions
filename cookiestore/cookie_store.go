package cookiestore

import (
	"encoding/json"
	"net/http"
	"time"
)

type CookieStore struct {
	*SecureCookie
}

func New(secret string) *CookieStore {
	return &CookieStore{
		NewSecureCookie([]byte(secret)),
	}
}

func (cs CookieStore) Get(ck *http.Cookie, pointer interface{}) error {
	if ck == nil || ck.Value == `` {
		return nil
	}
	b, err := cs.Decode(ck.Name, []byte(ck.Value), int64(ck.MaxAge))
	if err != nil || len(b) == 0 {
		return err
	}
	return json.Unmarshal(b, pointer)
}

func (cs *CookieStore) Save(res http.ResponseWriter, cookie *http.Cookie, data interface{}) error {
	encoded, err := cs.EncodeData(cookie.Name, data)
	if err != nil {
		return err
	}
	ck := *cookie // make a copy
	ck.Value = string(encoded)
	http.SetCookie(res, &ck)
	return nil
}

func (cs *CookieStore) Delete(res http.ResponseWriter, cookie *http.Cookie) {
	ck := *cookie // make a copy
	ck.Value = ""
	ck.MaxAge = -1
	ck.Expires = time.Unix(1, 0)
	http.SetCookie(res, &ck)
}

func (cs *CookieStore) EncodeData(cookieName string, data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return cs.Encode(cookieName, b)
}
