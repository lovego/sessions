package sessions

import (
	"fmt"

	"github.com/lovego/sessions/cookiestore"
)

func ExampleStore_assertCookieStore() {
	var cs interface{} = cookiestore.New("")
	if _, ok := cs.(Store); ok {
		fmt.Println("ok")
	}
	// Output: ok
}
