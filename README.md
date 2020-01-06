# gorequester
Golang chain style http request lib

# Install
```bash
go get github.com/cheetah-fun-gs/gorequester
```

# How
```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	requester "github.com/cheetah-fun-gs/gorequester"
)

func main() {
	// Get
	fmt.Println(requester.New("GET", "http://httpbin.org/get?abc=123").ReadString())
	fmt.Println(requester.New("GET", "http://httpbin.org/get?abc=123").AddRawQuery("def=456").ReadString())

	m := map[string]interface{}{
		"qewr": "424",
		"uio":  1231,
	}
	fmt.Println(requester.Get("http://httpbin.org/get?abc=123", m).ReadString())
	fmt.Println(requester.Get("http://httpbin.org/get?abc=123").AddRawQuery(m).ReadString())

	s := struct {
		FDI string `json:"fdi,omitempty"`
		IFU int    `json:"ifu,omitempty"`
	}{
		FDI: "9f0d",
		IFU: 413,
	}
	fmt.Println(requester.Get("http://httpbin.org/get?abc=123", s).ReadString())
	fmt.Println(requester.Get("http://httpbin.org/get?abc=123").AddRawQuery(s).ReadString())

	fmt.Println(requester.Get("http://httpbin.org/get?abc=123", m, s).ReadString())

	// Post Raw
	fmt.Println(requester.New("POST", "http://httpbin.org/post?abc=123").SetRawData(m).ReadString())
	fmt.Println(requester.Post("http://httpbin.org/post?abc=123").SetRawData(s).ReadString())
	fmt.Println(requester.PostData("http://httpbin.org/post?abc=123", "application/octet-stream", "abf").ReadString())
	fmt.Println(requester.PostJSON("http://httpbin.org/post?abc=123", s).ReadString())

	// Post Form
	m2 := map[string]interface{}{
		"qewr": "424",
		"uio":  []string{"dfa", "sf"},
		"fdf":  424,
		"gege": []int{3213, 4213},
	}
	fmt.Println(requester.PostForm("http://httpbin.org/post?abc=123", m2).ReadString())
	fmt.Println(requester.Post("http://httpbin.org/post?abc=123").SetFormFields(m2).ReadString())

	s2 := struct {
		FDI string   `json:"fdi,omitempty"`
		IFU int      `json:"ifu,omitempty"`
		IFI []int    `json:"ifi,omitempty"`
		GUE []string `json:"gue,omitempty"`
	}{
		FDI: "9f0d",
		IFU: 413,
		IFI: []int{12, 421},
		GUE: []string{"fs", "dsf"},
	}
	fmt.Println(requester.Post("http://httpbin.org/post?abc=123").SetFormFields(s2).ReadString())

	// Post File
	// f := &requester.FormFile{
	// 	FieldName: "file",
	// 	FileName:  "abc.txt",
	// 	FilePath:  "xxx",
	// }
	// fmt.Println(requester.Post("http://httpbin.org/post?abc=123").SetFormFields(s2).AddFormFile(f).ReadString())

	// JSON Resp
	r := map[string]interface{}{}
	if err := requester.Post("http://httpbin.org/post?abc=123").SetFormFields(s2).ReadJSON(&r); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	// Custom Client
	c := &http.Client{
		Timeout: time.Second * 5,
	}
	fmt.Println(requester.NewWithClient("GET", "http://httpbin.org/get?abc=123", c).ReadString())
}
```
