package sessions

import (
	"net/http"
	"time"
)

type CookieStore struct {
	*SecureCookie
}

func NewCookieStore(secret string) *CookieStore {
	return &CookieStore{
		NewSecureCookie([]byte(secret)),
	}
}

func (cs CookieStore) Get(req *http.Request, name string, maxAge int64) (*Session, error) {
	ck, _ := req.Cookie(name)
	if ck != nil && ck.Value != `` {
		data, err := cs.Decode(ck.Name, []byte(ck.Value), maxAge)
		if err != nil {
			return nil, err
		}
		return &Session{Cookie: ck, data: data}, nil
	}
	return nil, nil
}

func (cs *CookieStore) Save(res http.ResponseWriter, sess *Session) error {
	if err := cs.EncodeData(sess); err != nil {
		return err
	}
	http.SetCookie(res, sess.Cookie)
	return nil
}

func (cs *CookieStore) Delete(res http.ResponseWriter, sess *Session) {
	sess.Cookie.Value = ""
	sess.Cookie.MaxAge = -1
	sess.Cookie.Expires = time.Unix(1, 0)
	http.SetCookie(res, sess.Cookie)
}

func (cs *CookieStore) EncodeData(sess *Session) error {
	b, err := sess.jsonData()
	if err != nil {
		return err
	}
	encoded, err := cs.Encode(sess.Cookie.Name, b)
	if err != nil {
		return err
	}
	sess.Cookie.Value = string(encoded)
	return nil
}
