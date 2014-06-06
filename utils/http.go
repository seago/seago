package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"net/url"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 5,
	}
}

func PostJSON(url string, data interface{}) *simplejson.Json {
	logger.Info(url, data)
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		logger.Error(err)
		return nil
	}
	logger.Debug(string(dataBytes))
	request, _ := http.NewRequest("POST", url, bytes.NewReader(dataBytes))
	request.Header.Add("Content-Type", "application/json")
	httpClient := NewHttpClient()
	resp, err := httpClient.Do(request)
	if err != nil {
		logger.Error(err.Error())
		logger.Error("%#v", resp)
		return nil
	}
	return getJsonResponse(resp)
}

func PostForm(url string, data url.Values) *simplejson.Json {
	logger.Info(url, data)
	httpClient := NewHttpClient()
	resp, err := httpClient.PostForm(url, data)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return getJsonResponse(resp)
}

func GetUrl(url string, data url.Values) *simplejson.Json {
	logger.Info(url, data)
	httpClient := NewHttpClient()
	resp, err := httpClient.Get(fmt.Sprintf("%s?%s", url, data.Encode()))
	if err != nil {
		logger.Error(err)
		return nil
	}
	return getJsonResponse(resp)
}

func getJsonResponse(resp *http.Response) *simplejson.Json {
	logger.Debug(resp)
	var result []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(resp.Body)
		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)

			if err != nil && err != io.EOF {
				panic(err)
			}

			if n == 0 {
				break
			}
			result = append(result, buf...)
		}
	default:
		result, _ = ioutil.ReadAll(resp.Body)
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	logger.Debug("http resp(%d): %s", resp.StatusCode, string(result))
	jsonResponse, err := simplejson.NewJson(result)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return jsonResponse
}
