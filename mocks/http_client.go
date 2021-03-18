package mocks

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type MockHTTPClient struct {
	Resp  interface{}
	Error bool
}

func (m *MockHTTPClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	if m.Error {
		err = errors.New("Requested Error")
	}
	replyString, _ := json.Marshal(m.Resp)

	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(string(replyString))),
	}
	return
}
