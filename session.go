package sessions

import (
	"encoding/json"
	"net/http"
)

type Store interface {
	Get(req *http.Request, name string, maxAge int64) (*Session, error)
	Save(res http.ResponseWriter, sess *Session) error
	Delete(res http.ResponseWriter, sess *Session)
}

type Session struct {
	Cookie *http.Cookie
	data   []byte      // data got from request
	Data   interface{} // data to write to response
}

func (s *Session) GetData(pointer interface{}) error {
	if len(s.data) > 0 {
		return json.Unmarshal(s.data, pointer)
	}
	return nil
}

func (s *Session) jsonData() ([]byte, error) {
	return json.Marshal(s.Data)
}
