package gorequester

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	urlplus "github.com/cheetah-fun-gs/goplus/net/url"
)

func defaultClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 5,
	}
}

// Requester http.Request
type Requester struct {
	Error       error
	client      *http.Client
	url         *url.URL
	request     *http.Request
	method      string
	contentType string
	rawData     []byte
	formFields  url.Values
	formFiles   []*FormFile
}

// New 创建一个基础 Requester
func New(method, toURL string) *Requester {
	u, err := url.Parse(toURL)
	if err != nil {
		return &Requester{
			Error: err,
		}
	}
	return &Requester{
		client:     defaultClient(),
		url:        u,
		method:     method,
		rawData:    []byte{},
		formFields: map[string][]string{},
		formFiles:  []*FormFile{},
	}
}

// NewWithClient 创建一个基础 Requester
func NewWithClient(method, toURL string, client *http.Client) *Requester {
	u, err := url.Parse(toURL)
	if err != nil {
		return &Requester{
			Error: err,
		}
	}
	return &Requester{
		client:     client,
		url:        u,
		method:     method,
		rawData:    []byte{},
		formFields: map[string][]string{},
		formFiles:  []*FormFile{},
	}
}

// Client Client
func (req *Requester) Client() *http.Client {
	return req.client
}

// URL URL
func (req *Requester) URL() *url.URL {
	return req.url
}

// Request Request
func (req *Requester) Request() *http.Request {
	return req.request
}

// RawData RawData
func (req *Requester) RawData() []byte {
	return req.rawData
}

// FormFields FormFields
func (req *Requester) FormFields() url.Values {
	return req.formFields
}

// Get 创建一个 GET Requester
func Get(toURL string, v ...interface{}) *Requester {
	return New("GET", toURL).AddRawQuery(v...)
}

// Post 创建一个 POST Requester
func Post(toURL string) *Requester {
	return New("POST", toURL)
}

// PostData PostData
func PostData(toURL, contentType string, v interface{}) *Requester {
	req := New("POST", toURL).SetRawData(v)
	req.contentType = contentType
	return req
}

// PostJSON PostJSON
func PostJSON(toURL string, v interface{}) *Requester {
	req := New("POST", toURL).SetRawData(v)
	req.contentType = "application/json"
	return req
}

// PostForm PostForm
func PostForm(toURL string, v interface{}) *Requester {
	return New("POST", toURL).SetFormFields(v)
}

// v type in ( string, struct, map[string]string, map[string][]string,  map[string]int, map[string][]int )
func stringRawQuery(v ...interface{}) (string, error) {
	splits := []string{}
	for _, vv := range v {
		switch vv.(type) {
		case string:
			splits = append(splits, vv.(string))
		default:
			s, err := urlplus.ToRawQuery(vv)
			if err != nil {
				return "", err
			}
			splits = append(splits, s)
		}
	}
	return strings.Join(splits, "&"), nil
}

// AddRawQuery 追加 RawQuery
func (req *Requester) AddRawQuery(v ...interface{}) *Requester {
	if req.Error != nil {
		return req
	}
	s, err := stringRawQuery(v...)
	if err != nil {
		req.Error = err
		return req
	}
	if req.url.RawQuery == "" {
		req.url.RawQuery = s
	} else {
		req.url.RawQuery += "&" + s
	}
	return req
}

// SetRawQuery 重设 RawQuery
func (req *Requester) SetRawQuery(v ...interface{}) *Requester {
	if req.Error != nil {
		return req
	}
	s, err := stringRawQuery(v...)
	if err != nil {
		req.Error = err
		return req
	}
	req.url.RawQuery = s
	return req
}

// v type in ( string, []byte, struct, any json )
func byteRawData(v interface{}) ([]byte, error) {
	switch v.(type) {
	case string:
		return []byte(v.(string)), nil
	case []byte:
		return v.([]byte), nil
	default:
		d, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return d, nil
	}
}

// SetRawData 设置 RawData
func (req *Requester) SetRawData(v interface{}) *Requester {
	if req.Error != nil {
		return req
	}
	d, err := byteRawData(v)
	if err != nil {
		req.Error = err
		return req
	}
	req.rawData = d
	return req
}

// SetFormFields 设置 formField
// v type in ( string, struct, map[string]string, map[string][]string,  map[string]int, map[string][]int )
func (req *Requester) SetFormFields(v interface{}) *Requester {
	if req.Error != nil {
		return req
	}
	m, err := urlplus.ToValues(v)
	if err != nil {
		req.Error = err
		return req
	}
	req.formFields = m
	return req
}

// AddFormField 追加 formField
func (req *Requester) AddFormField(key string, val interface{}) *Requester {
	if req.Error != nil {
		return req
	}
	v, ok := req.formFields[key]
	if !ok {
		req.formFields[key] = []string{fmt.Sprintf("%v", val)}
	} else {
		v = append(v, fmt.Sprintf("%v", val))
	}
	return req
}

// SetFormField 设置 formField
func (req *Requester) SetFormField(key string, val interface{}) *Requester {
	if req.Error != nil {
		return req
	}
	req.formFields[key] = []string{fmt.Sprintf("%v", val)}
	return req
}

// AddFormFile 追加 formFile
func (req *Requester) AddFormFile(files ...*FormFile) *Requester {
	if req.Error != nil {
		return req
	}
	req.formFiles = append(req.formFiles, files...)
	return req
}

func (req *Requester) httpRequest() *Requester {
	if req.Error != nil {
		return req
	}
	if req.request != nil {
		return req
	}
	if req.method != "POST" && (len(req.rawData) != 0 || len(req.formFields) != 0 || len(req.formFiles) != 0) {
		req.Error = fmt.Errorf("have rawData, formField, formFile, use POST")
		return req
	}
	if len(req.rawData) != 0 && (len(req.formFields) != 0 || len(req.formFiles) != 0) {
		req.Error = fmt.Errorf("raw, form only choose one")
		return req
	}

	var body io.Reader
	if len(req.rawData) != 0 {
		body = strings.NewReader(string(req.rawData))
	} else {
		var contentType string
		var err error
		contentType, body, err = buildFormData(req.formFields, req.formFiles)
		if err != nil {
			req.Error = err
			return req
		}
		req.contentType = contentType
	}

	httpReq, err := http.NewRequest(req.method, req.url.String(), body)
	if err != nil {
		req.Error = err
		return req
	}
	if req.contentType != "" {
		httpReq.Header.Set("Content-Type", req.contentType)
	}
	req.request = httpReq
	return req
}

// SetHeader 设置请求头
func (req *Requester) SetHeader(key, val string) *Requester {
	req.httpRequest()
	if req.Error != nil {
		return req
	}
	req.request.Header.Set(key, val)
	return req
}

// AddHeader 添加请求头
func (req *Requester) AddHeader(key, val string) *Requester {
	req.httpRequest()
	if req.Error != nil {
		return req
	}
	req.request.Header.Add(key, val)
	return req
}

// Do 获取响应
func (req *Requester) Do() (*http.Response, error) {
	req.httpRequest()
	if req.Error != nil {
		return nil, req.Error
	}
	return req.client.Do(req.request)
}

// ReadData 获取响应二进制响应
func (req *Requester) ReadData() ([]byte, error) {
	resp, err := req.Do()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

// ReadString 获取响应字符串响应
func (req *Requester) ReadString() (string, error) {
	data, err := req.ReadData()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadJSON 获取响应JSON响应
func (req *Requester) ReadJSON(v interface{}) error {
	data, err := req.ReadData()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
