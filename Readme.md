# sessions
http sessions. now there is only cookie store, we will add redis store later.

[![Build Status](https://travis-ci.org/lovego/sessions.svg?branch=master)](https://travis-ci.org/lovego/sessions)
[![Coverage Status](https://img.shields.io/coveralls/github/lovego/sessions/master.svg)](https://coveralls.io/github/lovego/sessions?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/sessions?1)](https://goreportcard.com/report/github.com/lovego/sessions)
[![GoDoc](https://godoc.org/github.com/lovego/sessions?status.svg)](https://godoc.org/github.com/lovego/sessions)

## Usage
```go
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
```
