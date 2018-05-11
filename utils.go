package sinapay

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

//postForm POST请求
func postForm(urls string, postbody url.Values) ([]byte, error) {
	resp, err := http.PostForm(urls, postbody)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return handleBody(resp.Body)
}

func handleBody(respBody io.Reader) ([]byte, error) {
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//json2Map []byte JSON转map
func json2Map(b []byte) map[string]interface{} {
	var f interface{}
	json.Unmarshal(b, &f)
	return f.(map[string]interface{})
}
