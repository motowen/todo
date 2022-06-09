package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/app/router"
)

func HttpGet(path string, headers map[string]string) (resp *httptest.ResponseRecorder, err error) {
	resp, err = sendHttp("GET", path, "", headers)
	return
}

func HttpPost(path string, body string, headers map[string]string) (resp *httptest.ResponseRecorder, err error) {
	resp, err = sendHttp("POST", path, body, headers)
	return
}

func HttpPostAndMarshalBody(path string, body interface{}, headers map[string]string) (resp *httptest.ResponseRecorder, err error) {
	bodyInByte, err := json.Marshal(body)
	if err != nil {
		return
	}

	resp, err = sendHttp("POST", path, string(bodyInByte), headers)
	return
}

func HttpPatch(path string, body string, headers map[string]string) (resp *httptest.ResponseRecorder, err error) {
	resp, err = sendHttp("PATCH", path, body, headers)
	return
}

func HttpDelete(path string, body string, headers map[string]string) (resp *httptest.ResponseRecorder, err error) {
	resp, err = sendHttp("DELETE", path, body, headers)
	return
}

func sendHttp(method string, path string, body string, headers map[string]string) (resp *httptest.ResponseRecorder, err error) {
	resp = httptest.NewRecorder()

	req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	if err != nil {
		return
	}

	setHeader(req, headers)
	router.Router.ServeHTTP(resp, req)

	return
}

func setHeader(req *http.Request, headers map[string]string) {
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
