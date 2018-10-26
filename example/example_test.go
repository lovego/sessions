package example

import (
	"fmt"
	"net/http"

	"github.com/lovego/sessions"
)

var cookieStore = sessions.NewCookieStore("XXXX")
var cookie = &http.Cookie{
	Name:   "session",
	MaxAge: 86400 * 30,
}

type SessionData struct {
	UserId   int
	UserName string
}

func getSessionData(req *http.Request) *SessionData {
	sess, err := cookieStore.Get(req, cookie.Name, int64(cookie.MaxAge))
	if err != nil {
		fmt.Println("get session error:", err)
	}
	if sess == nil {
		return nil
	}
	var data SessionData
	if err := sess.GetData(&data); err != nil {
		fmt.Println("get session data error:", err)
	}
	return &data
}

func setSessionData(rw http.ResponseWriter, data *SessionData) {
	sess := &sessions.Session{
		Cookie: cookie,
		Data:   data,
	}

	err := cookieStore.Save(rw, sess)
	if err != nil {
		fmt.Println("set session error:", err)
	}
}

func deleteSession(rw http.ResponseWriter) {
	sess := &sessions.Session{
		Cookie: cookie,
	}

	cookieStore.Delete(rw, sess)
}
