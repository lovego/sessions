package example

import (
	"fmt"
	"net/http"

	"github.com/lovego/sessions/cookiestore"
)

var cookie = &http.Cookie{
	Name:   "session",
	MaxAge: 86400 * 30,
}
var cookieStore = cookiestore.New("XXXX")

type Session struct {
	UserId   int
	UserName string
}

func GetSession(req *http.Request) *Session {
	ck, err := req.Cookie(cookie.Name)
	if err != nil {
		fmt.Println("get cookie error:", err)
		return nil
	}
	if ck == nil || ck.Value == "" {
		return nil
	}
	ck.MaxAge = cookie.MaxAge

	var data Session
	if err := cookieStore.Get(ck, &data); err != nil {
		fmt.Println("get session error:", err)
		return nil
	}
	return &data
}

func SaveSession(rw http.ResponseWriter, data interface{}) {
	err := cookieStore.Save(rw, cookie, data)
	if err != nil {
		fmt.Println("save session error:", err)
	}
}

func DeleteSession(rw http.ResponseWriter) {
	cookieStore.Delete(rw, cookie)
}
