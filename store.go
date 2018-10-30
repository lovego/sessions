package sessions

import (
	"net/http"
)

type Store interface {
	Get(ck *http.Cookie, pointer interface{}) error
	Save(res http.ResponseWriter, ck *http.Cookie, data interface{}) error
	Delete(res http.ResponseWriter, ck *http.Cookie)
}
