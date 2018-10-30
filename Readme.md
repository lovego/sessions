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
```
