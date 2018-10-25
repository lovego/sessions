package sessions

import "fmt"

func ExampleSession_GetData() {
	s := Session{data: []byte(`{"UserId": 3}`)}
	var data struct {
		UserId int
	}
	if err := s.GetData(&data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", data)
	}
	// Output: {UserId:3}
}

func ExampleSession_GetData_nodata() {
	s := Session{}
	var data struct {
		UserId int
	}
	if err := s.GetData(&data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", data)
	}
	// Output: {UserId:0}
}
